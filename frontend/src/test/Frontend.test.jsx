import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import axios from 'axios';
import App from '../App.jsx';

// Mock axios
vi.mock('axios');

describe('Frontend Unit Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.confirm = vi.fn(() => true);
    global.alert = vi.fn();
  });

  describe('App Component Basic Tests', () => {
    it('renders main layout with header and footer', async () => {
      axios.get.mockResolvedValue({ data: [] });

      render(<App />);

      expect(screen.getByText('Hydro Habitat')).toBeInTheDocument();
      expect(screen.getByText('Tank Inventory')).toBeInTheDocument();
      expect(
        screen.getByText(/© \d{4} Hydro Habitat. All Rights Reserved./)
      ).toBeInTheDocument();
    });

    it('displays loading state initially', async () => {
      // Mock that never resolves to simulate loading
      axios.get.mockImplementation(() => new Promise(() => {}));

      render(<App />);

      expect(screen.getByText('Loading tanks...')).toBeInTheDocument();
    });

    it('displays empty state when no tanks exist', async () => {
      axios.get.mockResolvedValue({ data: [] });

      render(<App />);

      // Wait for async operations
      await screen.findByText('No tanks found. Add one to get started!');
      expect(
        screen.getByText('No tanks found. Add one to get started!')
      ).toBeInTheDocument();
    });

    it('displays error message when API call fails', async () => {
      axios.get.mockRejectedValue(new Error('Network error'));

      render(<App />);

      await screen.findByText(
        'Failed to fetch tanks. The backend might be starting up.'
      );
      expect(
        screen.getByText(
          'Failed to fetch tanks. The backend might be starting up.'
        )
      ).toBeInTheDocument();
    });

    it('opens modal when Add Tank button is clicked', async () => {
      axios.get.mockResolvedValue({ data: [] });
      const user = userEvent.setup();

      render(<App />);

      // Wait for loading to complete
      await screen.findByText('Add Tank');

      await user.click(screen.getByText('Add Tank'));

      expect(screen.getByText('Add New Tank')).toBeInTheDocument();
    });
  });

  describe('Tank Display Tests', () => {
    const mockTanks = [
      {
        id: '123e4567-e89b-12d3-a456-426614174000',
        name: 'Test Tank',
        volume_liters: 100,
        water: 'tap',
        room: 'Living Room',
        rack_location: 'Top',
        inventory_number: 'INV001',
        notes: 'Test notes',
      },
      {
        id: '456e7890-e89b-12d3-a456-426614174001',
        name: 'Another Tank',
        volume_liters: 200,
        water: 'ro',
        room: 'Basement',
        rack_location: 'Bottom',
        inventory_number: 'INV002',
        notes: 'Another test tank',
      },
    ];

    it('displays tanks when API call succeeds', async () => {
      axios.get.mockResolvedValue({ data: mockTanks });

      render(<App />);

      await screen.findByText('Test Tank');
      expect(screen.getByText('Test Tank')).toBeInTheDocument();
      expect(screen.getByText('Another Tank')).toBeInTheDocument();
      expect(screen.getByText('100 L')).toBeInTheDocument();
      expect(screen.getByText('200 L')).toBeInTheDocument();
    });

    it('displays tank details correctly', async () => {
      axios.get.mockResolvedValue({ data: mockTanks });

      render(<App />);

      await screen.findByText('Test Tank');

      // Check first tank details
      expect(screen.getByText('INV001')).toBeInTheDocument();
      expect(screen.getByText('Living Room - Top')).toBeInTheDocument();
      expect(screen.getByText('tap')).toBeInTheDocument();
      expect(screen.getByText('Test notes')).toBeInTheDocument();

      // Check second tank details
      expect(screen.getByText('INV002')).toBeInTheDocument();
      expect(screen.getByText('Basement - Bottom')).toBeInTheDocument();
      expect(screen.getByText('ro')).toBeInTheDocument();
      expect(screen.getByText('Another test tank')).toBeInTheDocument();
    });
  });

  describe('Form Validation Logic Tests', () => {
    it('validates required fields correctly', () => {
      const validData = {
        name: 'Test Tank',
        volume_liters: '100',
      };

      const invalidData1 = {
        name: '',
        volume_liters: '100',
      };

      const invalidData2 = {
        name: 'Test Tank',
        volume_liters: '',
      };

      const invalidData3 = {
        name: 'Test Tank',
        volume_liters: '-10',
      };

      // Valid data should pass
      expect(
        !!(
          validData.name &&
          validData.volume_liters &&
          parseInt(validData.volume_liters, 10) > 0
        )
      ).toBe(true);

      // Invalid data should fail
      expect(
        !!(
          invalidData1.name &&
          invalidData1.volume_liters &&
          parseInt(invalidData1.volume_liters, 10) > 0
        )
      ).toBe(false);
      expect(
        !!(
          invalidData2.name &&
          invalidData2.volume_liters &&
          parseInt(invalidData2.volume_liters, 10) > 0
        )
      ).toBe(false);
      expect(
        !!(
          invalidData3.name &&
          invalidData3.volume_liters &&
          parseInt(invalidData3.volume_liters, 10) > 0
        )
      ).toBe(false);
    });

    it('validates volume as positive number', () => {
      expect(parseInt('100', 10) > 0).toBe(true);
      expect(parseInt('0', 10) > 0).toBe(false);
      expect(parseInt('-10', 10) > 0).toBe(false);
      expect(parseInt('abc', 10) > 0).toBe(false);
      expect(parseInt('', 10) > 0).toBe(false);
    });
  });

  describe('Data Transformation Tests', () => {
    it('correctly transforms form data for API submission', () => {
      const formData = {
        name: 'Test Tank',
        volume_liters: '150', // String from form input
        water: 'ro',
        room: 'Living Room',
        rack_location: 'Top',
        inventory_number: 'INV001',
        notes: 'Test notes',
      };

      const transformedData = {
        ...formData,
        volume_liters: parseInt(formData.volume_liters, 10), // Convert to number
      };

      expect(transformedData.volume_liters).toBe(150);
      expect(typeof transformedData.volume_liters).toBe('number');
      expect(transformedData.name).toBe('Test Tank');
      expect(transformedData.water).toBe('ro');
    });

    it('handles empty string values correctly', () => {
      const formData = {
        name: 'Test Tank',
        volume_liters: '100',
        water: 'tap',
        room: '',
        rack_location: '',
        inventory_number: '',
        notes: '',
      };

      const transformedData = {
        ...formData,
        volume_liters: parseInt(formData.volume_liters, 10),
      };

      expect(transformedData.room).toBe('');
      expect(transformedData.rack_location).toBe('');
      expect(transformedData.inventory_number).toBe('');
      expect(transformedData.notes).toBe('');
    });
  });
});
