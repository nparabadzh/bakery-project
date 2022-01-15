import { combineReducers } from 'redux';

import userReducer from './user';
// import cartReducer from './cart/cartReducer';
// import directoryReducer from './directory/directoryReducer';
// import shopReducer from './shop/shopReducer';

const rootReducer = combineReducers({
  user: userReducer,
  // cart: cartReducer,
  // directory: directoryReducer,
  // shop: shopReducer,
});

export default rootReducer;
