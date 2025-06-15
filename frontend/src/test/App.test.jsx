import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import axios from 'axios';
import App from '../App.jsx';
import { mockTankList } from './test-utils.js';

// Mock axios
vi.mock('axios');

describe('App Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.confirm = vi.fn(() => true);
    global.alert = vi.fn();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

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
    axios.get.mockImplementation(() => new Promise(() => {})); // Never resolves

    render(<App />);

    expect(screen.getByText('Loading tanks...')).toBeInTheDocument();
  });

  it('displays tanks when API call succeeds', async () => {
    axios.get.mockResolvedValue({ data: mockTankList });

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText('Test Tank')).toBeInTheDocument();
    });
    expect(screen.getByText('Another Tank')).toBeInTheDocument();
  });

  it('displays error message when API call fails', async () => {
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

  it('displays empty state when no tanks exist', async () => {
    axios.get.mockResolvedValue({ data: [] });

    render(<App />);

    await waitFor(() => {
      expect(
        screen.getByText('No tanks found. Add one to get started!')
      ).toBeInTheDocument();
    });
  });

  it('opens modal when Add Tank button is clicked', async () => {
    axios.get.mockResolvedValue({ data: [] });
    const user = userEvent.setup();

    render(<App />);

    await waitFor(() => {
      expect(screen.getByText('Add Tank')).toBeInTheDocument();
    });

    await user.click(screen.getByText('Add Tank'));

    expect(screen.getByText('Add New Tank')).toBeInTheDocument();
  });
});
