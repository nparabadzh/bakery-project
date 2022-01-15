import React from 'react';
import './App.css';
import { Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import Home from './pages/home';
import Cakes from './pages/cakes';
import Pies from './pages/pies';
import Pastries from './pages/pastries';
import Other from './pages/other';
import SignIn from './components/SignIn';
import SignUp from './components/SignUp';

function App() {
  return (
    <div>
      <div id="bg" />
      <Navbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="cakes" element={<Cakes />} />
        <Route path="pies" element={<Pies />} />
        <Route path="pastries" element={<Pastries />} />
        <Route path="other" element={<Other />} />
        <Route path="signin" element={<SignIn />} />
        <Route path="signup" element={<SignUp />} />
      </Routes>
    </div>
  );
}

export default App;
