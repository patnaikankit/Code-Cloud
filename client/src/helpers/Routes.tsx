import React from 'react';
import { Route, Routes } from 'react-router-dom';
import Home from '../pages/Home';
import Editor from '../pages/Editor';

const RoutesContainer: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/:id" element={<Editor />} />
    </Routes>
  );
}

export default RoutesContainer;
