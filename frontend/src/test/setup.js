import '@testing-library/jest-dom';
import { vi } from 'vitest';

// Mock window.confirm and window.alert for tests
global.confirm = vi.fn(() => true);
global.alert = vi.fn();

// Mock axios for tests
vi.mock('axios', () => ({
  default: {
    get: vi.fn(() => Promise.resolve({ data: [] })),
    post: vi.fn(() => Promise.resolve({ data: {} })),
    put: vi.fn(() => Promise.resolve({ data: {} })),
    delete: vi.fn(() => Promise.resolve({ data: {} })),
  },
}));
