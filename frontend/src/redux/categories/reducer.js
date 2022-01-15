import { CategoriesActionTypes } from './types';

const INITIAL_STATE = {
  categories: {},
};

const categoriesReducer = (state = INITIAL_STATE, action) => {
  switch (action.type) {
    case CategoriesActionTypes.SET_CATEGORIES:
      return {
        ...state,
        categories: action.payload,
      };

    default:
      return state;
  }
};

export default categoriesReducer;
