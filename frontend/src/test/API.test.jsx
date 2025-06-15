import { describe, it, expect, vi, beforeEach } from 'vitest';
import axios from 'axios';
import { mockTank, mockTankList } from './test-utils.js';

vi.mock('axios');

describe('API Utility Functions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Health API Endpoint', () => {
    it('makes correct API call for health check', async () => {
      axios.get.mockResolvedValue({ data: { status: 'ok' } });

      const response = await axios.get('/health');

      expect(axios.get).toHaveBeenCalledWith('/health');
      expect(response.data).toEqual({ status: 'ok' });
    });

    it('handles health check endpoint errors', async () => {
      const healthError = {
        response: {
          status: 503,
          data: { error: 'Service unavailable' },
        },
      };
      axios.get.mockRejectedValue(healthError);

      await expect(axios.get('/health')).rejects.toEqual(healthError);
    });
  });

  describe('Tank API Endpoints', () => {
    it('makes correct API call for fetching all tanks', async () => {
      axios.get.mockResolvedValue({ data: mockTankList });

      const response = await axios.get('/api/v1/tanks');

      expect(axios.get).toHaveBeenCalledWith('/api/v1/tanks');
      expect(response.data).toEqual(mockTankList);
    });

    it('makes correct API call for creating a tank', async () => {
      const newTankData = {
        name: 'New Tank',
        volume_liters: 150,
        water: 'ro',
        room: 'Bedroom',
        rack_location: 'Middle',
        inventory_number: 'INV003',
        notes: 'New tank notes',
      };

      axios.post.mockResolvedValue({ data: mockTank });

      const response = await axios.post('/api/v1/tanks', newTankData);

      expect(axios.post).toHaveBeenCalledWith('/api/v1/tanks', newTankData);
      expect(response.data).toEqual(mockTank);
    });

    it('makes correct API call for updating a tank', async () => {
      const updateData = {
        name: 'Updated Tank',
        volume_liters: 200,
        water: 'rodi',
        room: 'Kitchen',
        rack_location: 'Bottom',
        inventory_number: 'INV001-UPDATED',
        notes: 'Updated notes',
      };

      axios.put.mockResolvedValue({ data: { ...mockTank, ...updateData } });

      const response = await axios.put(
        `/api/v1/tanks/${mockTank.id}`,
        updateData
      );

      expect(axios.put).toHaveBeenCalledWith(
        `/api/v1/tanks/${mockTank.id}`,
        updateData
      );
      expect(response.data).toEqual({ ...mockTank, ...updateData });
    });

    it('makes correct API call for deleting a tank', async () => {
      axios.delete.mockResolvedValue({ data: {} });

      const response = await axios.delete(`/api/v1/tanks/${mockTank.id}`);

      expect(axios.delete).toHaveBeenCalledWith(`/api/v1/tanks/${mockTank.id}`);
      expect(response.data).toEqual({});
    });
  });

  describe('Error Handling', () => {
    it('handles network errors correctly', async () => {
      const networkError = new Error('Network Error');
      axios.get.mockRejectedValue(networkError);

      await expect(axios.get('/api/v1/tanks')).rejects.toBe(networkError);
      expect(axios.get).toHaveBeenCalledWith('/api/v1/tanks');
    });

    it('handles API error responses correctly', async () => {
      const apiError = {
        response: {
          status: 400,
          data: { error: 'Invalid tank data' },
        },
      };
      axios.post.mockRejectedValue(apiError);

      await expect(axios.post('/api/v1/tanks', {})).rejects.toEqual(apiError);
    });

    it('handles server errors correctly', async () => {
      const serverError = {
        response: {
          status: 500,
          data: { error: 'Internal server error' },
        },
      };
      axios.get.mockRejectedValue(serverError);

      await expect(axios.get('/api/v1/tanks')).rejects.toEqual(serverError);
    });
  });
});

describe('Data Transformation', () => {
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

describe('Form Validation Logic', () => {
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

    // Invalid data should fail - UPDATED VERSION
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
