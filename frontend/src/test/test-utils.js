import { render } from '@testing-library/react';
import { vi } from 'vitest';

// Mock tank data for testing
export const mockTank = {
  id: '123e4567-e89b-12d3-a456-426614174000',
  name: 'Test Tank',
  volume_liters: 100,
  water: 'tap',
  room: 'Living Room',
  rack_location: 'Top',
  inventory_number: 'INV001',
  notes: 'Test notes',
  created_at: '2023-01-01T00:00:00Z',
  updated_at: '2023-01-01T00:00:00Z',
};

export const mockTankList = [
  mockTank,
  {
    id: '456e7890-e89b-12d3-a456-426614174001',
    name: 'Another Tank',
    volume_liters: 200,
    water: 'ro',
    room: 'Basement',
    rack_location: 'Bottom',
    inventory_number: 'INV002',
    notes: 'Another test tank',
    created_at: '2023-01-02T00:00:00Z',
    updated_at: '2023-01-02T00:00:00Z',
  },
];

// Custom render function that includes common providers if needed
export const renderWithProviders = (ui, options = {}) => {
  return render(ui, {
    ...options,
  });
};

// Helper to create mock axios responses
export const createMockAxios = (responses = {}) => {
  const defaultResponses = {
    get: Promise.resolve({ data: mockTankList }),
    post: Promise.resolve({ data: mockTank }),
    put: Promise.resolve({ data: mockTank }),
    delete: Promise.resolve({ data: {} }),
  };

  return {
    get: vi.fn(() => responses.get || defaultResponses.get),
    post: vi.fn(() => responses.post || defaultResponses.post),
    put: vi.fn(() => responses.put || defaultResponses.put),
    delete: vi.fn(() => responses.delete || defaultResponses.delete),
  };
};

// Helper to wait for async operations
export const waitForAsync = () =>
  new Promise((resolve) => setTimeout(resolve, 0));
