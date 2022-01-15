import { CategoriesActionTypes } from './types';

export const setCategories = (categories) => ({
  type: CategoriesActionTypes.SET_CATEGORIES,
  payload: categories,
});
