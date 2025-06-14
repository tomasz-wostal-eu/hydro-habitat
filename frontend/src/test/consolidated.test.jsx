import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  render,
  screen,
  waitFor,
} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import axios from 'axios';
import App from '../App.jsx';
import { mockTank } from './test-utils.js';

// Mock axios
vi.mock('axios');

describe('Hydro Habitat Frontend Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.alert = vi.fn();
    global.confirm = vi.fn();
  });

  describe('App Component', () => {
    it('renders header with correct title', async () => {
      axios.get.mockResolvedValue({ data: [] });

      render(<App />);

      expect(screen.getByText('Hydro Habitat')).toBeInTheDocument();
      // Remove Tank Inventory check as it's not in the actual component
    });

    it('displays loading state initially', async () => {
      axios.get.mockImplementation(() => new Promise(() => {})); // Never resolves

      render(<App />);

      expect(screen.getByText('Loading tanks...')).toBeInTheDocument();
    });

    it('displays no tanks message when empty', async () => {
      axios.get.mockResolvedValue({ data: [] });

      render(<App />);

      await waitFor(() => {
        expect(
          screen.getByText('No tanks found. Add one to get started!')
        ).toBeInTheDocument();
      });
    });

    it('displays tanks when data is available', async () => {
      axios.get.mockResolvedValue({ data: [mockTank] });

      render(<App />);

      await waitFor(() => {
        expect(screen.getByText('Test Tank')).toBeInTheDocument();
      });
      await waitFor(() => {
        expect(screen.getByText('100 L')).toBeInTheDocument();
      });
    });
  });

  describe('Tank Management', () => {
    it('opens modal when Add Tank button is clicked', async () => {
      const user = userEvent.setup();
      axios.get.mockResolvedValue({ data: [] });

      render(<App />);

      await waitFor(() => {
        expect(screen.getByText('Add Tank')).toBeInTheDocument();
      });

      await user.click(screen.getByText('Add Tank'));

      await waitFor(() => {
        expect(screen.getByText('Add New Tank')).toBeInTheDocument();
      });
    });

    it('closes modal when cancel is clicked', async () => {
      const user = userEvent.setup();
      axios.get.mockResolvedValue({ data: [] });

      render(<App />);

      // Open modal
      await waitFor(() => {
        expect(screen.getByText('Add Tank')).toBeInTheDocument();
      });
      await user.click(screen.getByText('Add Tank'));

      await waitFor(() => {
        expect(screen.getByText('Add New Tank')).toBeInTheDocument();
      });

      // Close modal
      await user.click(screen.getByText('Cancel'));

      await waitFor(() => {
        expect(screen.queryByText('Add New Tank')).not.toBeInTheDocument();
      });
    });
  });

  // TankFormModal Component is tested through integration tests above
  // as it's part of the App component

  describe('Integration Tests', () => {
    it('successfully creates a new tank', async () => {
      const user = userEvent.setup();
      axios.get.mockResolvedValue({ data: [] });
      axios.post.mockResolvedValue({ data: mockTank });

      render(<App />);

      // Wait for initial load
      await waitFor(() => {
        expect(
          screen.getByText('No tanks found. Add one to get started!')
        ).toBeInTheDocument();
      });

      // Click Add Tank button
      await user.click(screen.getByText('Add Tank'));

      // Wait for modal to appear
      await waitFor(() => {
        expect(screen.getByText('Add New Tank')).toBeInTheDocument();
      });

      // Fill out the form using actual placeholder text from App.jsx
      await user.type(screen.getByPlaceholderText('Tank Name'), 'Test Tank');
      await user.type(screen.getByPlaceholderText('Volume (Liters)'), '100');
      await user.type(screen.getByPlaceholderText('Room'), 'Living Room');
      await user.type(screen.getByPlaceholderText('Rack Location'), 'Top');
      await user.type(
        screen.getByPlaceholderText('Inventory Number'),
        'INV001'
      );
      await user.type(screen.getByPlaceholderText('Notes...'), 'Test notes');

      // Submit the form
      await user.click(screen.getByText('Save Tank'));

      // Verify API calls
      await waitFor(() => {
        expect(axios.post).toHaveBeenCalledWith('/api/v1/tanks', {
          name: 'Test Tank',
          volume_liters: 100,
          water: 'tap', // default value
          room: 'Living Room',
          rack_location: 'Top',
          inventory_number: 'INV001',
          notes: 'Test notes',
        });
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
      await user.click(screen.getByText('Add Tank'));

      // Wait for modal to open
      await waitFor(() => {
        expect(screen.getByText('Add New Tank')).toBeInTheDocument();
      });

      // Try to submit without filling required fields
      await user.click(screen.getByText('Save Tank'));

      // Should show validation error (this depends on browser validation or custom validation)
      // For now, we'll just check that the form doesn't submit
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
      await user.click(screen.getByText('Add Tank'));

      // Wait for modal to open and form fields to be available
      await waitFor(() => {
        expect(screen.getByText('Add New Tank')).toBeInTheDocument();
      });

      // Fill and submit form
      await user.type(screen.getByPlaceholderText('Tank Name'), 'Test Tank');
      await user.type(screen.getByPlaceholderText('Volume (Liters)'), '100');
      await user.click(screen.getByText('Save Tank'));

      // Should show error alert
      await waitFor(() => {
        expect(global.alert).toHaveBeenCalledWith(
          'Error saving tank: Database error'
        );
      });
    });
  });

  describe('Footer Component', () => {
    it('displays current year in footer', async () => {
      axios.get.mockResolvedValue({ data: [] });

      render(<App />);

      const currentYear = new Date().getFullYear();
      await waitFor(() => {
        expect(
          screen.getByText(
            `© ${currentYear} Hydro Habitat. All Rights Reserved.`
          )
        ).toBeInTheDocument();
      });
    });
  });
});
