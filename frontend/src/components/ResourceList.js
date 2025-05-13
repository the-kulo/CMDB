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
      setError(null); // 重置错误状态
      try {
        let endpoint = '/api/vms'; // Initialize with default

        // 根据 resourceType 选择 API 端点
        switch (resourceType) {
          case 'database':
            endpoint = '/api/databases';
            break;
          case 'sqldatabase':
            endpoint = '/api/sqldatabase';
            break;
          case 'mysqlflexible':
            endpoint = '/api/mysqlflexible';
            break;
          case 'sqlserver':
            endpoint = '/api/sqlserver';
            break;
          default: 
            endpoint = '/api/vms'; 
            break;
        }

        const response = await fetch(`${endpoint}`);
        if (!response.ok) {
          const errorData = await response.text(); // 获取错误详情
          throw new Error(`请求失败: ${response.status} - ${errorData}`);
        }
        const data = await response.json();
        setResources(data || []); // 确保 data 不是 null 或 undefined
        setLoading(false);
      } catch (err) {
        console.error("获取资源时出错:", err); // 在控制台打印详细错误
        setError(err.message);
        setLoading(false);
      }
    };

    fetchResources();
  }, [resourceType]); // 依赖项是 resourceType

  if (loading) {
    return <div className="loading">loading...</div>;
  }

  if (error) {
    return <div className="error">加载资源时出错: {error}</div>;
  }

  const renderTableHeaders = () => {
    // 根据资源类型渲染不同的表头
    switch (resourceType) {
      case 'database': // 为 'database' 类型定义表头
        return (
          <tr>
            <th>Name</th>
            <th>Database ID</th>
            <th>Type</th>
            <th>Server</th>
            <th>Location</th>
            <th>Status</th>
            <th>Owner</th>
          </tr>
        );
      // 保留其他特定类型的表头，如果需要的话
      case 'sqldatabase':
        // ... (如果 App.js 未映射，则保留)
        return (
          <tr>
            <th>Name</th>
            <th>Server</th>
            <th>Location</th>
            <th>Status</th>
            <th>Owner</th>
          </tr>
        );
      case 'mysqlflexible':
        // ... (如果 App.js 未映射，则保留)
        return (
          <tr>
            <th>Name</th>
            <th>Version</th>
            <th>Location</th>
            <th>Status</th>
            <th>Owner</th>
          </tr>
        );
      case 'sqlserver':
        // ... (如果 App.js 未映射，则保留)
        return (
          <tr>
            <th>Name</th>
            <th>Version</th>
            <th>Location</th>
            <th>Status</th>
            <th>Owner</th>
          </tr>
        );
      default: // VM
        return (
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Location</th>
            <th>Status</th>
            <th>Owner</th>
          </tr>
        );
    }
  };

  const renderTableRows = () => {
    if (!resources || resources.length === 0) {
      return <tr><td colSpan="7">没有找到资源</td></tr>; // 调整 colSpan
    }

    return resources.map((resource) => {
      // 根据资源类型渲染不同的表格行
      switch (resourceType) {
        case 'database': // 为 'database' 类型渲染行
          return (
            <tr key={resource.database_id || resource.id}> {/* 使用 database_id 作为 key */}
              <td>{resource.name || 'N/A'}</td>
              <td>{resource.database_id || 'N/A'}</td>
              <td>{resource.db_type || 'N/A'}</td>
              <td>{resource.server || 'N/A'}</td>
              <td>{resource.location || 'N/A'}</td>
              <td>{resource.status || 'N/A'}</td>
              <td>{resource.owner || 'N/A'}</td>
            </tr>
          );
        // 保留其他特定类型的行渲染，如果需要的话
        case 'sqldatabase':
          // ... (如果 App.js 未映射，则保留)
          return (
            <tr key={resource.database_id || resource.id}>
              <td>{resource.name || 'N/A'}</td>
              <td>{resource.server || 'N/A'}</td>
              <td>{resource.location || 'N/A'}</td>
              <td>{resource.status || 'N/A'}</td>
              <td>{resource.owner || 'N/A'}</td>
            </tr>
          );
        case 'mysqlflexible':
          // ... (如果 App.js 未映射，则保留)
          return (
            <tr key={resource.database_id || resource.id}>
              <td>{resource.name || 'N/A'}</td>
              <td>{resource.version || 'N/A'}</td>
              <td>{resource.location || 'N/A'}</td>
              <td>{resource.status || 'N/A'}</td>
              <td>{resource.owner || 'N/A'}</td>
            </tr>
          );
        case 'sqlserver':
          // ... (如果 App.js 未映射，则保留)
          return (
            <tr key={resource.database_id || resource.id}>
              <td>{resource.name || 'N/A'}</td>
              <td>{resource.version || 'N/A'}</td>
              <td>{resource.location || 'N/A'}</td>
              <td>{resource.status || 'N/A'}</td>
              <td>{resource.owner || 'N/A'}</td>
            </tr>
          );
        default: // VM
          return (
            <tr key={resource.vm_id || resource.id}> {/* 使用 vm_id 作为 key */}
              <td>{resource.name || 'N/A'}</td>
              <td>{resource.type || 'N/A'}</td>
              <td>{resource.location || 'N/A'}</td>
              <td>{resource.status || 'N/A'}</td>
              <td>{resource.owner || 'N/A'}</td>
            </tr>
          );
      }
    });
  };

  return (
    <div className="resource-list-container">
      <div className="resource-list">
        <table className="resource-table">
          <thead>
            {renderTableHeaders()}
          </thead>
          <tbody>
            {renderTableRows()}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default ResourceList;