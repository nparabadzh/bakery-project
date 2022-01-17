import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { useDispatch } from 'react-redux';

import { clearItemFromCart, addItem, removeItem } from '../redux/cart/actions';

const useStyles = makeStyles(() => ({
  checkoutItem: {
    width: '100%',
    display: 'flex',
    minHeight: 100,
    borderBottom: '1px solid darkgrey',
    padding: '15px 0',
    fontSize: 20,
    alignItems: 'center',
  },
  imageContainer: {
    width: '23%',
    paddingRight: 15,
  },
  name: {
    width: '23%',
  },
  quantity: {
    width: '23%',
    display: 'flex',
  },
  arrow: {
    cursor: 'pointer',
  },
  value: {
    margin: '0 10px',
  },
  price: {
    width: '23%',
  },
  removeButton: {
    paddingLeft: 12,
    cursor: 'pointer',
  },
}));

const CheckoutItem = ({ cartItem }) => {
  const dispatch = useDispatch();
  const classes = useStyles();
  const { name, photo, price, quantity } = cartItem;
  return (
    <div className={classes.checkoutItem}>
      <div className={classes.imageContainer}>
        <img style={{ width: '100%', height: 'auto' }} src={photo} alt="item" />
      </div>
      <span className={classes.name}>{name}</span>
      <span className={classes.quantity}>
        <div
          className={classes.arror}
          onClick={() => dispatch(removeItem(cartItem))}
        >
          &#10094;
        </div>
        <span className={classes.value}>{quantity}</span>
        <div
          className={classes.arror}
          onClick={() => dispatch(addItem(cartItem))}
        >
          &#10095;
        </div>
      </span>
      <span className={classes.price}>{price} bgn</span>
      <div
        className={classes.removeButton}
        onClick={() => dispatch(clearItemFromCart(cartItem))}
      >
        &#10005;
      </div>
    </div>
  );
};

export default CheckoutItem;
