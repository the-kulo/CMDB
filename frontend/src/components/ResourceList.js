import React, { useState, useEffect } from 'react';
import '../styles/ResourceList.css';
import 'boxicons/css/boxicons.min.css';

const ResourceList = ({ resourceType = 'vm' }) => {
  const [resources, setResources] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchResources = async () => {
      setLoading(true);
      try {
        let endpoint = '/api/vms';
        
        // 根据资源类型选择不同的API端点
        switch(resourceType) {
          case 'sqldatabase':
            endpoint = '/api/sqldatabases';
            break;
          case 'mysqlflexible':
            endpoint = '/api/mysqlflexible';
            break;
          case 'sqlserver':
            endpoint = '/api/sqlservers';
            break;
          case 'database': // 保留原有的数据库类型（向后兼容）
            endpoint = '/api/databases';
            break;
          default:
            endpoint = '/api/vms';
        }
        
        const response = await fetch(`http://localhost:8080${endpoint}`);
        if (!response.ok) {
          throw new Error(`请求失败: ${response.status}`);
        }
        const data = await response.json();
        setResources(data);
        setLoading(false);
      } catch (err) {
        setError(err.message);
        setLoading(false);
      }
    };

    fetchResources();
  }, [resourceType]);

  if (loading) {
    return <div className="loading">loading...</div>;
  }

  if (error) {
    return <div className="error">error: {error}</div>;
  }

  const renderTableHeaders = () => {
    // 根据资源类型渲染不同的表头
    switch(resourceType) {
      case 'sqldatabase':
        return (
          <tr>
            <th>Name</th>
            <th>Sub ID</th>
            <th>Location</th>
            <th>Resource Group</th>
            <th>Status</th>
            <th>IT Owner</th>
          </tr>
        );
      case 'mysqlflexible':
        return (
          <tr>
            <th>Name</th>
            <th>Sub ID</th>
            <th>Location</th>
            <th>Version</th>
            <th>Status</th>
            <th>IT Owner</th>
          </tr>
        );
      case 'sqlserver':
        return (
          <tr>
            <th>Name</th>
            <th>Sub ID</th>
            <th>Location</th>
            <th>Version</th>
            <th>Status</th>
            <th>IT Owner</th>
          </tr>
        );
      case 'database': // 保留原有的数据库类型（向后兼容）
        return (
          <tr>
            <th>Name</th>
            <th>Sub ID</th>
            <th>Location</th>
            <th>IT Owner</th>
          </tr>
        );
      default:
        return (
          <tr>
            <th>Name</th>
            <th>Sub ID</th>
            <th>Location</th>
            <th>OS Type</th>
            <th>IT Owner</th>
          </tr>
        );
    }
  };

  const renderTableRows = () => {
    return resources.map((resource) => {
      // 根据资源类型渲染不同的表格行
      switch(resourceType) {
        case 'sqldatabase':
          return (
            <tr key={resource.id}>
              <td>{resource.name}</td>
              <td>{resource.id}</td>
              <td>{resource.location}</td>
              <td>{resource.server || '未指定'}</td>
              <td>{resource.status || '未知'}</td>
              <td>{resource.owner || '未指定'}</td>
            </tr>
          );
        case 'mysqlflexible':
          return (
            <tr key={resource.id}>
              <td>{resource.name}</td>
              <td>{resource.id}</td>
              <td>{resource.location}</td>
              <td>{resource.version || '未指定'}</td>
              <td>{resource.status || '未知'}</td>
              <td>{resource.owner || '未指定'}</td>
            </tr>
          );
        case 'sqlserver':
          return (
            <tr key={resource.id}>
              <td>{resource.name}</td>
              <td>{resource.id}</td>
              <td>{resource.location}</td>
              <td>{resource.version || '未指定'}</td>
              <td>{resource.status || '未知'}</td>
              <td>{resource.owner || '未指定'}</td>
            </tr>
          );
        case 'database': // 保留原有的数据库类型（向后兼容）
          return (
            <tr key={resource.id}>
              <td>{resource.name}</td>
              <td>{resource.id}</td>
              <td>{resource.location}</td>
              <td>{resource.owner || '未指定'}</td>
            </tr>
          );
        default:
          return (
            <tr key={resource.id}>
              <td>{resource.name}</td>
              <td>{resource.id}</td>
              <td>{resource.location}</td>
              <td>{resource.size || '未指定'}</td>
              <td>{resource.owner || '未指定'}</td>
            </tr>
          );
      }
    });
  };

  return (
    <div className="resource-list-container">
      <div className="resource-list">
        {resources.length === 0 ? (
          <p>没有找到资源</p>
        ) : (
          <table className="resource-table">
            <thead>
              {renderTableHeaders()}
            </thead>
            <tbody>
              {renderTableRows()}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default ResourceList;