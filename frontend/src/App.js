import React from 'react';
import ResourceList from './components/ResourceList';
import './styles/App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>CMDB</h1>
        <p>资源配置管理数据库</p>
      </header>
      <main>
        <ResourceList />
      </main>
      <footer>
        <p>&copy; {new Date().getFullYear()} CMDB</p>
      </footer>
    </div>
  );
}

export default App;