import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import axios from 'axios';
import App from '../App.jsx';
import { mockTank, mockTankList } from './test-utils.js';

vi.mock('axios');

describe('Tank Management Integration Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.confirm = vi.fn(() => true);
    global.alert = vi.fn();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

  describe('Tank Creation Flow', () => {
    it('successfully creates a new tank', async () => {
      const user = userEvent.setup();

      // Mock initial empty state and successful creation
      axios.get
        .mockResolvedValueOnce({ data: [] }) // Initial load
        .mockResolvedValueOnce({ data: [mockTank] }); // After creation

      axios.post.mockResolvedValue({ data: mockTank });

      render(<App />);

      // Wait for initial load
      await waitFor(() => {
        expect(
          screen.getByText('No tanks found. Add one to get started!')
        ).toBeInTheDocument();
      });

      // Click Add Tank button
      await act(async () => {
        await user.click(screen.getByText('Add Tank'));
      });

      // Fill out the form
      await act(async () => {
        await user.type(screen.getByPlaceholderText('Tank Name'), 'Test Tank');
        await user.type(screen.getByPlaceholderText('Volume (Liters)'), '100');
        await user.selectOptions(screen.getByDisplayValue('Tap Water'), 'tap');
        await user.type(screen.getByPlaceholderText('Room'), 'Living Room');
        await user.type(screen.getByPlaceholderText('Rack Location'), 'Top');
        await user.type(
          screen.getByPlaceholderText('Inventory Number'),
          'INV001'
        );
        await user.type(screen.getByPlaceholderText('Notes...'), 'Test notes');
      });

      // Submit the form
      await act(async () => {
        await user.click(screen.getByText('Save Tank'));
      });

      // Verify API calls
      expect(axios.post).toHaveBeenCalledWith('/api/v1/tanks', {
        name: 'Test Tank',
        volume_liters: 100,
        water: 'tap',
        room: 'Living Room',
        rack_location: 'Top',
        inventory_number: 'INV001',
        notes: 'Test notes',
      });

      // Verify tanks list is refreshed
      await waitFor(() => {
        expect(axios.get).toHaveBeenCalledTimes(2);
      });
    });

    it('shows validation error for missing required fields', async () => {
      const user = userEvent.setup();
      axios.get.mockResolvedValue({ data: [] });

      render(<App />);

      await waitFor(() => {
        expect(screen.getByText('Add Tank')).toBeInTheDocument();
      });

      // Click Add Tank button
      await act(async () => {
        await user.click(screen.getByText('Add Tank'));
      });

      // Try to submit without filling required fields
      await act(async () => {
        await user.click(screen.getByText('Save Tank'));
      });

      // Should show validation error
      expect(global.alert).toHaveBeenCalledWith(
        'Tank Name and a valid positive Volume are required.'
      );
      expect(axios.post).not.toHaveBeenCalled();
    });

    it('handles API error during creation', async () => {
      const user = userEvent.setup();
      axios.get.mockResolvedValue({ data: [] });
      axios.post.mockRejectedValue({
        response: { data: { error: 'Database error' } },
      });

      render(<App />);

      await waitFor(() => {
        expect(screen.getByText('Add Tank')).toBeInTheDocument();
      });

      // Click Add Tank and wait for modal to appear
      await act(async () => {
        await user.click(screen.getByText('Add Tank'));
      });

      // Wait for modal to open and form fields to be available
      await waitFor(() => {
        expect(screen.getByPlaceholderText('Tank Name')).toBeInTheDocument();
      });

      // Fill and submit form
      await act(async () => {
        await user.type(screen.getByPlaceholderText('Tank Name'), 'Test Tank');
        await user.type(screen.getByPlaceholderText('Volume (Liters)'), '100');
        await user.click(screen.getByText('Save Tank'));
      });

      // Should show error alert
      await waitFor(() => {
        expect(global.alert).toHaveBeenCalledWith(
          'Error saving tank: Database error'
        );
      });
    });
  });

  describe('Tank Update Flow', () => {
    it('successfully updates an existing tank', async () => {
      const user = userEvent.setup();

      axios.get
        .mockResolvedValueOnce({ data: [mockTank] }) // Initial load
        .mockResolvedValueOnce({
          data: [{ ...mockTank, name: 'Updated Tank' }],
        }); // After update

      axios.put.mockResolvedValue({
        data: { ...mockTank, name: 'Updated Tank' },
      });

      render(<App />);

      // Wait for initial load
      await waitFor(() => {
        expect(screen.getByText('Test Tank')).toBeInTheDocument();
      });

      // Click edit button (we need to find it by role since it's an icon)
      const editButtons = screen.getAllByRole('button');
      const editButton = editButtons.find(
        (button) =>
          button.innerHTML.includes('svg') &&
          button.className.includes('text-blue')
      );
      await act(async () => {
        await user.click(editButton);
      });

      // Verify form is pre-filled
      expect(screen.getByDisplayValue('Test Tank')).toBeInTheDocument();
      expect(screen.getByDisplayValue('100')).toBeInTheDocument();

      // Update the tank name
      const nameInput = screen.getByDisplayValue('Test Tank');
      await act(async () => {
        await user.clear(nameInput);
        await user.type(nameInput, 'Updated Tank');
      });

      // Submit the form
      await act(async () => {
        await user.click(screen.getByText('Save Tank'));
      });

      // Verify API call
      expect(axios.put).toHaveBeenCalledWith(`/api/v1/tanks/${mockTank.id}`, {
        name: 'Updated Tank',
        volume_liters: 100,
        water: 'tap',
        room: 'Living Room',
        rack_location: 'Top',
        inventory_number: 'INV001',
        notes: 'Test notes',
      });
    });
  });

  describe('Tank Deletion Flow', () => {
    it('successfully deletes a tank', async () => {
      const user = userEvent.setup();

      axios.get
        .mockResolvedValueOnce({ data: [mockTank] }) // Initial load
        .mockResolvedValueOnce({ data: [] }); // After deletion

      axios.delete.mockResolvedValue({ data: {} });

      render(<App />);

      // Wait for initial load
      await waitFor(() => {
        expect(screen.getByText('Test Tank')).toBeInTheDocument();
      });

      // Click delete button
      const deleteButtons = screen.getAllByRole('button');
      const deleteButton = deleteButtons.find(
        (button) =>
          button.innerHTML.includes('svg') &&
          button.className.includes('text-red')
      );
      await act(async () => {
        await user.click(deleteButton);
      });

      // Verify confirmation and API call
      expect(global.confirm).toHaveBeenCalledWith(
        'Are you sure you want to delete this tank?'
      );
      expect(axios.delete).toHaveBeenCalledWith(`/api/v1/tanks/${mockTank.id}`);
    });

    it('cancels deletion when user declines confirmation', async () => {
      const user = userEvent.setup();
      global.confirm.mockReturnValue(false); // User cancels

      axios.get.mockResolvedValue({ data: [mockTank] });

      render(<App />);

      await waitFor(() => {
        expect(screen.getByText('Test Tank')).toBeInTheDocument();
      });

      // Click delete button
      const deleteButtons = screen.getAllByRole('button');
      const deleteButton = deleteButtons.find(
        (button) =>
          button.innerHTML.includes('svg') &&
          button.className.includes('text-red')
      );
      await act(async () => {
        await user.click(deleteButton);
      });

      // Verify confirmation was shown but API was not called
      expect(global.confirm).toHaveBeenCalled();
      expect(axios.delete).not.toHaveBeenCalled();
    });

    it('handles API error during deletion', async () => {
      const user = userEvent.setup();

      axios.get.mockResolvedValue({ data: [mockTank] });
      axios.delete.mockRejectedValue({
        response: { data: { error: 'Cannot delete tank' } },
      });

      render(<App />);

      await waitFor(() => {
        expect(screen.getByText('Test Tank')).toBeInTheDocument();
      });

      // Click delete button
      const deleteButtons = screen.getAllByRole('button');
      const deleteButton = deleteButtons.find(
        (button) =>
          button.innerHTML.includes('svg') &&
          button.className.includes('text-red')
      );
      await act(async () => {
        await user.click(deleteButton);
      });

      // Should show error alert
      await waitFor(() => {
        expect(global.alert).toHaveBeenCalledWith(
          'Error deleting tank: Cannot delete tank'
        );
      });
    });
  });

  describe('Tank List Display', () => {
    it('displays multiple tanks correctly', async () => {
      axios.get.mockResolvedValue({ data: mockTankList });

      render(<App />);

      await waitFor(() => {
        expect(screen.getByText('Test Tank')).toBeInTheDocument();
      });
      await waitFor(() => {
        expect(screen.getByText('Another Tank')).toBeInTheDocument();
      });

      // Verify tank details are displayed
      expect(screen.getByText('100 L')).toBeInTheDocument();
      expect(screen.getByText('200 L')).toBeInTheDocument();
      expect(screen.getByText('Living Room - Top')).toBeInTheDocument();
      expect(screen.getByText('Basement - Bottom')).toBeInTheDocument();
    });

    it('handles network errors gracefully', async () => {
      axios.get.mockRejectedValue(new Error('Network error'));

      render(<App />);

      await waitFor(() => {
        expect(
          screen.getByText(
            'Failed to fetch tanks. The backend might be starting up.'
          )
        ).toBeInTheDocument();
      });
    });
  });
});
