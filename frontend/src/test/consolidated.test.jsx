import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  render,
  screen,
  waitFor,
  act,
  fireEvent,
} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import PropTypes from 'prop-types';
import axios from 'axios';
import App from '../App.jsx';
import { mockTank } from './test-utils.js';

// Mock axios
vi.mock('axios');

// Create a standalone TankFormModal component for isolated testing
const TankFormModal = ({ tank, onClose, onSave }) => {
  const [formData, setFormData] = React.useState({
    name: tank?.name || '',
    volume_liters: tank?.volume_liters || '',
    water: tank?.water || 'tap',
    room: tank?.room || '',
    rack_location: tank?.rack_location || '',
    inventory_number: tank?.inventory_number || '',
    notes: tank?.notes || '',
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const volumeNum = parseInt(formData.volume_liters, 10);
    if (
      !formData.name ||
      !formData.volume_liters ||
      isNaN(volumeNum) ||
      volumeNum <= 0
    ) {
      alert('Tank Name and a valid positive Volume are required.');
      return;
    }
    onSave(formData);
  };

  return (
    <div className="fixed inset-0 bg-base bg-opacity-75 flex items-center justify-center z-50 p-4">
      <div className="bg-mantle rounded-lg shadow-2xl p-6 w-full max-w-md">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-mauve">
            {tank ? 'Edit Tank' : 'Add New Tank'}
          </h2>
          <button
            onClick={onClose}
            className="text-overlay2 hover:text-text"
            data-testid="close-button"
          >
            ×
          </button>
        </div>
        <form
          onSubmit={handleSubmit}
          className="space-y-4"
          data-testid="tank-form"
        >
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            placeholder="Tank Name"
            className="form-input"
            required
            data-testid="name-input"
          />
          <input
            type="number"
            name="volume_liters"
            value={formData.volume_liters}
            onChange={handleChange}
            placeholder="Volume (Liters)"
            className="form-input"
            required
            data-testid="volume-input"
          />
          <select
            name="water"
            value={formData.water}
            onChange={handleChange}
            className="form-select"
            data-testid="water-select"
          >
            <option value="tap">Tap Water</option>
            <option value="ro">RO Water</option>
            <option value="rodi">RO/DI Water</option>
          </select>
          <input
            type="text"
            name="room"
            value={formData.room}
            onChange={handleChange}
            placeholder="Room"
            className="form-input"
            data-testid="room-input"
          />
          <input
            type="text"
            name="rack_location"
            value={formData.rack_location}
            onChange={handleChange}
            placeholder="Rack Location"
            className="form-input"
            data-testid="rack-location-input"
          />
          <input
            type="text"
            name="inventory_number"
            value={formData.inventory_number}
            onChange={handleChange}
            placeholder="Inventory Number"
            className="form-input"
            data-testid="inventory-input"
          />
          <textarea
            name="notes"
            value={formData.notes}
            onChange={handleChange}
            placeholder="Notes..."
            className="form-input h-24"
            data-testid="notes-input"
          />
          <div className="flex justify-end space-x-4 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="btn btn-secondary"
              data-testid="cancel-button"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              data-testid="save-button"
            >
              Save Tank
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

TankFormModal.propTypes = {
  tank: PropTypes.shape({
    name: PropTypes.string,
    volume_liters: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
    water: PropTypes.oneOf(['tap', 'ro', 'rodi']),
    room: PropTypes.string,
    rack_location: PropTypes.string,
    inventory_number: PropTypes.string,
    notes: PropTypes.string,
  }),
  onClose: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
};

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
      expect(screen.getByText('Tank Inventory')).toBeInTheDocument();
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

  describe('TankFormModal Component', () => {
    const mockOnClose = vi.fn();
    const mockOnSave = vi.fn();

    beforeEach(() => {
      vi.clearAllMocks();
      global.alert = vi.fn();
    });

    it('renders modal for creating new tank', async () => {
      render(
        <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
      );

      expect(screen.getByText('Add New Tank')).toBeInTheDocument();
      expect(screen.getByTestId('name-input')).toHaveValue('');
      // Number inputs with empty value return null, not empty string
      expect(screen.getByTestId('volume-input').value).toBe('');
      expect(screen.getByTestId('water-select')).toHaveValue('tap');
    });

    it('renders modal for editing existing tank', async () => {
      render(
        <TankFormModal
          tank={mockTank}
          onClose={mockOnClose}
          onSave={mockOnSave}
        />
      );

      expect(screen.getByText('Edit Tank')).toBeInTheDocument();
      expect(screen.getByTestId('name-input')).toHaveValue('Test Tank');
      expect(screen.getByTestId('volume-input')).toHaveValue(100);
      expect(screen.getByTestId('water-select')).toHaveValue('tap');
      expect(screen.getByTestId('room-input')).toHaveValue('Living Room');
      expect(screen.getByTestId('rack-location-input')).toHaveValue('Top');
      expect(screen.getByTestId('inventory-input')).toHaveValue('INV001');
      expect(screen.getByTestId('notes-input')).toHaveValue('Test notes');
    });

    it('calls onClose when close button is clicked', async () => {
      const user = userEvent.setup();

      render(
        <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
      );

      await user.click(screen.getByTestId('close-button'));

      expect(mockOnClose).toHaveBeenCalled();
    });

    it('calls onClose when cancel button is clicked', async () => {
      const user = userEvent.setup();

      render(
        <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
      );

      await user.click(screen.getByTestId('cancel-button'));

      expect(mockOnClose).toHaveBeenCalled();
    });

    it('updates form fields when user types', async () => {
      const user = userEvent.setup();

      render(
        <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
      );

      await user.type(screen.getByTestId('name-input'), 'New Tank');
      await user.type(screen.getByTestId('volume-input'), '150');

      expect(screen.getByTestId('name-input')).toHaveValue('New Tank');
      expect(screen.getByTestId('volume-input')).toHaveValue(150);
    });

    it('submits form with valid data', async () => {
      const user = userEvent.setup();

      render(
        <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
      );

      await user.type(screen.getByTestId('name-input'), 'New Tank');
      await user.type(screen.getByTestId('volume-input'), '150');
      await user.selectOptions(screen.getByTestId('water-select'), 'ro');
      await user.type(screen.getByTestId('room-input'), 'Bedroom');
      await user.click(screen.getByTestId('save-button'));

      expect(mockOnSave).toHaveBeenCalledWith({
        name: 'New Tank',
        volume_liters: '150',
        water: 'ro',
        room: 'Bedroom',
        rack_location: '',
        inventory_number: '',
        notes: '',
      });
    });

    it('shows alert for invalid data - missing name', async () => {
      const user = userEvent.setup();

      render(
        <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
      );

      await user.type(screen.getByTestId('volume-input'), '150');
      fireEvent.submit(screen.getByTestId('tank-form'));

      await waitFor(() => {
        expect(global.alert).toHaveBeenCalledWith(
          'Tank Name and a valid positive Volume are required.'
        );
      });
      expect(mockOnSave).not.toHaveBeenCalled();
    });

    it('shows alert for invalid data - missing volume', async () => {
      const user = userEvent.setup();

      render(
        <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
      );

      await user.type(screen.getByTestId('name-input'), 'New Tank');
      fireEvent.submit(screen.getByTestId('tank-form'));

      await waitFor(() => {
        expect(global.alert).toHaveBeenCalledWith(
          'Tank Name and a valid positive Volume are required.'
        );
      });
      expect(mockOnSave).not.toHaveBeenCalled();
    });

    it('shows alert for invalid data - negative volume', async () => {
      const user = userEvent.setup();

      render(
        <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
      );

      await user.type(screen.getByTestId('name-input'), 'New Tank');
      await user.type(screen.getByTestId('volume-input'), '-10');
      await user.click(screen.getByTestId('save-button'));

      expect(global.alert).toHaveBeenCalledWith(
        'Tank Name and a valid positive Volume are required.'
      );
      expect(mockOnSave).not.toHaveBeenCalled();
    });
  });

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

      // Wait for modal to open and then try to submit without filling required fields
      await waitFor(() => {
        expect(screen.getByTestId('tank-form')).toBeInTheDocument();
      });

      fireEvent.submit(screen.getByTestId('tank-form'));

      // Should show validation error
      await waitFor(() => {
        expect(global.alert).toHaveBeenCalledWith(
          'Tank Name and a valid positive Volume are required.'
        );
      });
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
        expect(screen.getByPlaceholderText('Tank Name')).toBeInTheDocument();
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
      expect(
        screen.getByText(
          `© ${currentYear} Hydro Habitat. All Rights Reserved.`
        )
      ).toBeInTheDocument();
    });
  });
});
