import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(() => ({
  cartItem: {
    width: '100',
    display: 'flex',
    height: 80,
    marginBottom: 15,
  },
  itemDetails: {
    width: '70%',
    color: 'black',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'flex-start',
    justifyContent: 'center',
    padding: '10px 20px',
  },
  name: {
    fontSize: 16,
  },
}));

const CartItem = ({ item: { photo, price, name, quantity } }) => {
  const classes = useStyles();
  debugger;
  return (
    <div className={classes.cartItem}>
      <img style={{ width: '30%' }} src={photo} alt="item" />
      <div className={classes.itemDetails}>
        <span className={classes.name}>{name}</span>
        <span className={classes.price}>
          {quantity} x {price} bgn
        </span>
      </div>
    </div>
  );
};

export default CartItem;
