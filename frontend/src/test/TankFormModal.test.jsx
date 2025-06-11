import React from 'react';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import PropTypes from 'prop-types';

// Create a standalone TankFormModal component for testing
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
    const volumeNum = parseFloat(formData.volume_liters);
    if (
      !formData.name.trim() ||
      !formData.volume_liters ||
      formData.volume_liters === '' ||
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
            X
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
            data-testid="name-input"
          />
          <input
            type="number"
            name="volume_liters"
            value={formData.volume_liters}
            onChange={handleChange}
            placeholder="Volume (Liters)"
            className="form-input"
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
            data-testid="inventory-number-input"
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

describe('TankFormModal Component', () => {
  const mockOnClose = vi.fn();
  const mockOnSave = vi.fn();
  const mockTank = {
    id: 1,
    name: 'Test Tank',
    volume_liters: 100,
    water: 'tap',
    room: 'Living Room',
  };

  beforeEach(() => {
    vi.clearAllMocks();
    global.alert = vi.fn();
  });

  it('renders modal for creating new tank', () => {
    render(
      <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
    );

    expect(screen.getByText('Add New Tank')).toBeInTheDocument();
    expect(screen.getByTestId('name-input')).toHaveValue('');
    expect(screen.getByTestId('volume-input')).toHaveValue(null);
    expect(screen.getByTestId('water-select')).toHaveValue('tap');
  });

  it('renders modal for editing existing tank', () => {
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
  });

  it('calls onClose when close button is clicked', async () => {
    const user = userEvent.setup();
    render(
      <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
    );

    await act(async () => {
      await user.click(screen.getByTestId('close-button'));
    });

    expect(mockOnClose).toHaveBeenCalled();
  });

  it('calls onClose when cancel button is clicked', async () => {
    const user = userEvent.setup();
    render(
      <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
    );

    await act(async () => {
      await user.click(screen.getByTestId('cancel-button'));
    });

    expect(mockOnClose).toHaveBeenCalled();
  });

  it('updates form fields when user types', async () => {
    const user = userEvent.setup();
    render(
      <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
    );

    await act(async () => {
      await user.clear(screen.getByTestId('name-input'));
      await user.type(screen.getByTestId('name-input'), 'New Tank');
      // For number inputs, type event might result in a number value or null if empty
      // Clearing and then typing '150'
      await user.clear(screen.getByTestId('volume-input'));
      await user.type(screen.getByTestId('volume-input'), '150');
    });

    expect(screen.getByTestId('name-input')).toHaveValue('New Tank');
    expect(screen.getByTestId('volume-input')).toHaveValue(150);
  });

  it('submits form with valid data', async () => {
    const user = userEvent.setup();
    render(
      <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
    );

    await act(async () => {
      await user.type(screen.getByTestId('name-input'), 'My New Tank');
      await user.clear(screen.getByTestId('volume-input'));
      await user.type(screen.getByTestId('volume-input'), '150');
      await user.selectOptions(screen.getByTestId('water-select'), 'ro');
      await user.type(screen.getByTestId('room-input'), 'Bedroom');

      await user.click(screen.getByTestId('save-button'));
    });

    expect(mockOnSave).toHaveBeenCalledWith({
      name: 'My New Tank',
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

    await act(async () => {
      await user.clear(screen.getByTestId('name-input'));
      await user.clear(screen.getByTestId('volume-input'));
      await user.type(screen.getByTestId('volume-input'), '150');
    });

    await act(async () => {
      await user.click(screen.getByTestId('save-button'));
    });

    expect(global.alert).toHaveBeenCalledWith(
      'Tank Name and a valid positive Volume are required.'
    );
    expect(mockOnSave).not.toHaveBeenCalled();
  });

  it('shows alert for invalid data - missing volume', async () => {
    const user = userEvent.setup();
    render(
      <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
    );

    await act(async () => {
      await user.clear(screen.getByTestId('volume-input'));
      await user.clear(screen.getByTestId('name-input'));
      await user.type(screen.getByTestId('name-input'), 'Test Tank');
    });

    await act(async () => {
      await user.click(screen.getByTestId('save-button'));
    });

    expect(global.alert).toHaveBeenCalledWith(
      'Tank Name and a valid positive Volume are required.'
    );
    expect(mockOnSave).not.toHaveBeenCalled();
  });

  it('shows alert for invalid data - negative volume', async () => {
    const user = userEvent.setup();
    render(
      <TankFormModal tank={null} onClose={mockOnClose} onSave={mockOnSave} />
    );

    await act(async () => {
      await user.type(screen.getByTestId('name-input'), 'New Tank');
      await user.type(screen.getByTestId('volume-input'), '-10');
    });

    await act(async () => {
      await user.click(screen.getByTestId('save-button'));
    });

    expect(global.alert).toHaveBeenCalledWith(
      'Tank Name and a valid positive Volume are required.'
    );
    expect(mockOnSave).not.toHaveBeenCalled();
  });
});
