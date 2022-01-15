import { combineReducers } from 'redux';

import userReducer from './user/reducer';
import categoriesReducer from './categories/reducer';
// import directoryReducer from './directory/directoryReducer';
// import shopReducer from './shop/shopReducer';

const rootReducer = combineReducers({
  user: userReducer,
  categories: categoriesReducer,
  // cart: cartReducer,
  // directory: directoryReducer,
  // shop: shopReducer,
});

export default rootReducer;
