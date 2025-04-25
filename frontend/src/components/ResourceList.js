import React, { useState, useEffect } from 'react';
import '../styles/ResourceList.css';

const ResourceList = () => {
  const [resources, setResources] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchResources = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/resources');
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
  }, []);

  if (loading) {
    return <div className="loading">加载中...</div>;
  }

  if (error) {
    return <div className="error">错误: {error}</div>;
  }

  return (
    <div className="resource-list-container">
      <h1>Azure 资源列表</h1>
      <div className="resource-list">
        {resources.length === 0 ? (
          <p>没有找到资源</p>
        ) : (
          resources.map((resource) => (
            <div className="resource-card" key={resource.id}>
              <h2>{resource.name}</h2>
              <div className="resource-details">
                <p><strong>资源ID:</strong> {resource.id}</p>
                <p><strong>位置:</strong> {resource.location}</p>
                <p><strong>所有者:</strong> {resource.owner || '未指定'}</p>
                <p><strong>类型:</strong> {resource.type}</p>
              </div>
              <div className="resource-tags">
                <h3>标签:</h3>
                {Object.keys(resource.tags).length > 0 ? (
                  <ul>
                    {Object.entries(resource.tags).map(([key, value]) => (
                      <li key={key}>
                        <strong>{key}:</strong> {value}
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p>无标签</p>
                )}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default ResourceList;