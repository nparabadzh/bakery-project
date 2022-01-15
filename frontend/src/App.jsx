import React, { useEffect } from 'react';
import './App.css';
import { Routes, Route } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import Navbar from './components/Navbar';
import Home from './pages/home';
import Cakes from './pages/cakes';
import Pies from './pages/pies';
import Pastries from './pages/pastries';
import Other from './pages/other';
import SignIn from './pages/sign/SignIn';
import SignUp from './pages/sign/SignUp';
import axios from 'axios';
import { setCategories } from './redux/categories/actions';

function App() {
  const dispatch = useDispatch();
  useEffect(() => {
    axios
      .get('/categories')
      .then((resp) => {
        const categories = resp.data.reduce((accumulator, value) => {
          accumulator[value.id] = value.category_name;
          return accumulator;
        }, {});
        dispatch(setCategories(categories));
      })
      .catch((error) => console.log(`Unable to get categories ${error}`));
  }, []);

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
