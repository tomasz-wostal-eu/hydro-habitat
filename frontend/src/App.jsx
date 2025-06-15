import { useState, useEffect, useCallback } from 'react';
import PropTypes from 'prop-types';
import axios from 'axios';
import { Plus, Edit, Trash2, Droplets, X } from 'lucide-react';

// Main Application Component
function App() {
  return (
    <div className="min-h-screen bg-crust font-sans">
      <Header />
      <main className="p-4 sm:p-6 lg:p-8">
        <TankManager />
      </main>
      <Footer />
    </div>
  );
}

const Header = () => (
  <header className="bg-mantle shadow-md">
    <div className="container mx-auto px-4 py-3 flex justify-between items-center">
      <div className="flex items-center space-x-3">
        <Droplets className="text-blue h-8 w-8" />
        <h1 className="text-2xl font-bold text-text">Hydro Habitat</h1>
      </div>
    </div>
  </header>
);

const TankManager = () => {
  const [tanks, setTanks] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTank, setEditingTank] = useState(null);

  const API_URL = '/api/v1/tanks';

  const fetchTanks = useCallback(async () => {
    try {
      setIsLoading(true);
      const response = await axios.get(API_URL);
      setTanks(response.data || []);
      setError(null);
    } catch {
      setError('Failed to fetch tanks. The backend might be starting up.');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchTanks();
  }, [fetchTanks]);

  const handleOpenModal = (tank = null) => {
    setEditingTank(tank);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingTank(null);
  };

  const handleSaveTank = async (tankData) => {
    const dataToSend = {
      ...tankData,
      volume_liters: parseInt(tankData.volume_liters, 10),
    };

    try {
      if (editingTank) {
        await axios.put(`${API_URL}/${editingTank.id}`, dataToSend);
      } else {
        await axios.post(API_URL, dataToSend);
      }
      fetchTanks();
      handleCloseModal();
    } catch (error) {
      alert('Error saving tank: ' + (error.response?.data?.error || error.message));
    }
  };

  const handleDeleteTank = async (tankId) => {
    if (window.confirm('Are you sure you want to delete this tank?')) {
      try {
        await axios.delete(`${API_URL}/${tankId}`);
        fetchTanks();
      } catch (err) {
        alert(
          'Error deleting tank: ' + (err.response?.data?.error || err.message)
        );
      }
    }
  };

  return (
    <div className="container mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-3xl font-semibold text-subtext1">Tank Inventory</h2>
        <button
          onClick={() => handleOpenModal()}
          className="btn btn-primary flex items-center space-x-2"
        >
          <Plus size={20} />
          <span>Add Tank</span>
        </button>
      </div>

      {isLoading && (
        <p className="text-center text-overlay1">Loading tanks...</p>
      )}
      {error && <p className="text-center text-red">{error}</p>}

      {!isLoading && !error && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {tanks.length > 0 ? (
            tanks.map((tank) => (
              <TankCard
                key={tank.id}
                tank={tank}
                onEdit={handleOpenModal}
                onDelete={handleDeleteTank}
              />
            ))
          ) : (
            <p className="text-center text-overlay1 col-span-full">
              No tanks found. Add one to get started!
            </p>
          )}
        </div>
      )}

      {isModalOpen && (
        <TankFormModal
          tank={editingTank}
          onClose={handleCloseModal}
          onSave={handleSaveTank}
        />
      )}
    </div>
  );
};

const TankCard = ({ tank, onEdit, onDelete }) => (
  <div className="bg-mantle rounded-lg shadow-lg p-5 flex flex-col justify-between transition-transform transform hover:-translate-y-1">
    <div>
      <h3 className="text-xl font-bold text-peach truncate">{tank.name}</h3>
      <p className="text-sm text-overlay2 mb-3">
        {tank.inventory_number || 'No ID'}
      </p>
      <div className="space-y-2 text-sm text-subtext1">
        <p>
          <strong>Volume:</strong> {tank.volume_liters} L
        </p>
        <p>
          <strong>Location:</strong> {tank.room || 'N/A'} -{' '}
          {tank.rack_location || 'N/A'}
        </p>
        <p>
          <strong>Water:</strong>{' '}
          <span
            className="font-semibold capitalize"
            style={{
              color:
                tank.water === 'rodi'
                  ? '#8aadf4'
                  : tank.water === 'ro'
                    ? '#91d7e3'
                    : '#a6da95',
            }}
          >
            {tank.water}
          </span>
        </p>
        <p className="text-xs text-overlay1 pt-2 italic">
          {tank.notes || 'No additional notes.'}
        </p>
      </div>
    </div>
    <div className="flex justify-end space-x-2 mt-4">
      <button
        onClick={() => onEdit(tank)}
        className="p-2 text-blue hover:text-sapphire"
      >
        <Edit size={18} />
      </button>
      <button
        onClick={() => onDelete(tank.id)}
        className="p-2 text-red hover:text-maroon"
      >
        <Trash2 size={18} />
      </button>
    </div>
  </div>
);

TankCard.propTypes = {
  tank: PropTypes.shape({
    id: PropTypes.oneOfType([PropTypes.string, PropTypes.number]).isRequired,
    name: PropTypes.string.isRequired,
    inventory_number: PropTypes.string,
    volume_liters: PropTypes.oneOfType([PropTypes.string, PropTypes.number])
      .isRequired,
    room: PropTypes.string,
    rack_location: PropTypes.string,
    water: PropTypes.oneOf(['tap', 'ro', 'rodi']).isRequired,
    notes: PropTypes.string,
  }).isRequired,
  onEdit: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
};

const TankFormModal = ({ tank, onClose, onSave }) => {
  const [formData, setFormData] = useState({
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
            <X size={24} />
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
            data-testid="inventory-input"
          />
          <textarea
            name="notes"
            value={formData.notes}
            onChange={handleChange}
            placeholder="Notes..."
            className="form-input h-24"
            data-testid="notes-input"
          ></textarea>
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

const Footer = () => (
  <footer className="bg-mantle mt-8" role="contentinfo">
    <div className="container mx-auto py-4 text-center text-sm text-overlay1">
      <p>
        &copy; {new Date().getFullYear()} Hydro Habitat. All Rights Reserved.
      </p>
    </div>
  </footer>
);

export default App;
