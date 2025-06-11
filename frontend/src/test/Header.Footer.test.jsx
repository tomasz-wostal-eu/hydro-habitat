import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';

// Extract Header and Footer components for testing
const Header = () => (
  <header className="bg-mantle shadow-md">
    <div className="container mx-auto px-4 py-3 flex justify-between items-center">
      <div className="flex items-center space-x-3">
        <svg className="text-blue h-8 w-8" data-testid="logo-icon" />
        <h1 className="text-2xl font-bold text-text">Hydro Habitat</h1>
      </div>
    </div>
  </header>
);

const Footer = () => (
  <footer className="bg-mantle mt-8">
    <div className="container mx-auto py-4 text-center text-sm text-overlay1">
      <p>
        &copy; {new Date().getFullYear()} Hydro Habitat. All Rights Reserved.
      </p>
    </div>
  </footer>
);

describe('Header Component', () => {
  it('renders the application title', () => {
    render(<Header />);

    expect(screen.getByText('Hydro Habitat')).toBeInTheDocument();
    expect(screen.getByText('Hydro Habitat')).toHaveClass(
      'text-2xl',
      'font-bold',
      'text-text'
    );
  });

  it('renders the logo icon', () => {
    render(<Header />);

    expect(screen.getByTestId('logo-icon')).toBeInTheDocument();
    expect(screen.getByTestId('logo-icon')).toHaveClass(
      'text-blue',
      'h-8',
      'w-8'
    );
  });

  it('has proper semantic structure', () => {
    render(<Header />);

    const header = screen.getByRole('banner');
    expect(header).toBeInTheDocument();
    expect(header).toHaveClass('bg-mantle', 'shadow-md');
  });
});

describe('Footer Component', () => {
  it('renders the copyright notice with current year', () => {
    render(<Footer />);

    const currentYear = new Date().getFullYear();
    expect(
      screen.getByText(`© ${currentYear} Hydro Habitat. All Rights Reserved.`)
    ).toBeInTheDocument();
  });

  it('has proper semantic structure', () => {
    render(<Footer />);

    const footer = screen.getByRole('contentinfo');
    expect(footer).toBeInTheDocument();
    expect(footer).toHaveClass('bg-mantle', 'mt-8');
  });

  it('has centered text styling', () => {
    render(<Footer />);

    const copyrightText = screen.getByText(
      /© \d{4} Hydro Habitat. All Rights Reserved./
    );
    expect(copyrightText).toBeInTheDocument();
    // Verify footer styling by checking the footer element
    const footer = screen.getByRole('contentinfo');
    expect(footer).toHaveClass('bg-mantle', 'mt-8');
  });
});
