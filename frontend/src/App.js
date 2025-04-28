import React, { useState, useRef, useEffect } from 'react';
import ResourceList from './components/ResourceList';
import './styles/App.css';

function App() {
  const [activeTab, setActiveTab] = useState('vm');
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const toggleDropdown = () => {
    setIsOpen(!isOpen);
  };

  const selectOption = (option) => {
    setActiveTab(option.toLowerCase());
    // setIsOpen(false);
  };

  const renderContent = () => {
    return <ResourceList resourceType={activeTab} />;
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1 style={{ paddingLeft:'20px'}}>CMDB</h1>
      </header>
      <div className="App-container">
        <nav className="App-sidebar">
          <div className="select-menu" ref={dropdownRef}>
            <div className="select-btn" onClick={toggleDropdown}>
              <span>Resource List</span>
              <i className={`bx ${isOpen ? 'bx-chevron-up' : 'bx-chevron-down'}`}></i>
            </div>
            <ul className={`options ${isOpen ? 'active' : ''}`}>
              <li className={activeTab === 'vm' ? 'selected' : ''} onClick={() => selectOption('vm')}>VM</li>
              <li className={activeTab === 'sqldatabase' ? 'selected' : ''} onClick={() => selectOption('sqldatabase')}>SQL Database</li>
              <li className={activeTab === 'sqlserver' ? 'selected' : ''} onClick={() => selectOption('sqlserver')}>SQL Server</li>
              <li className={activeTab === 'mysqlflexible' ? 'selected' : ''} onClick={() => selectOption('mysqlflexible')}>MySQL Flexible Server</li>
            </ul>
          </div>
        </nav>
        <main className="App-content">
          <div className="search-box">
            <input type="text" placeholder="Have a nice day ^-^" />
            <button>Search</button>
          </div>
          {renderContent()}
        </main>
      </div>
    </div>
  );
}

export default App;