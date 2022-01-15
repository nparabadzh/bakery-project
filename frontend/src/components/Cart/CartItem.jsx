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

const CartItem = ({ item: { imageUrl, price, name, quantity } }) => {
  const classes = useStyles();
  console.log(classes);
  return (
    <div className={classes.cartItem}>
      <img style={{ width: '30%' }} src={imageUrl} alt="item" />
      <div className={classes.itemDetails}>
        <span className={classes.name}>{name}</span>
        <span className={classes.price}>
          {quantity} x ${price}
        </span>
      </div>
    </div>
  );
};

export default CartItem;
