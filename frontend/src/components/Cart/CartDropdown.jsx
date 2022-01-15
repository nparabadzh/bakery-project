import React from 'react';

import CustomButton from '../CustomButton';
import CartItem from './CartItem';

import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(() => ({
  cartDropdown: {
    position: 'absolute',
    width: 240,
    height: 340,
    display: 'flex',
    flexDirection: 'column',
    padding: 20,
    border: '1px solid black',
    backgroundColor: 'white',
    top: 90,
    right: 40,
    zIndex: 5,
  },
  emptyMsg: {
    fontSize: 18,
    margin: '50px auto',
  },
  cartItems: {
    height: 240,
    display: 'flex',
    flexDirection: 'column',
    overflow: 'scroll',
  },
  button: {
    marginTop: 'auto',
  },
}));

const CartDropdown = ({ cartItems }) => {
  const classes = useStyles();
  return (
    <div className={classes.cartDropdown}>
      <div className={classes.cartItems}>
        {cartItems.length ? (
          cartItems.map((cartItem) => (
            <CartItem key={cartItem.id} item={cartItem} />
          ))
        ) : (
          <span className={classes.emptyMsg}>Your cart is empty</span>
        )}
      </div>
      <div className={classes.button}>
        <CustomButton
          onClick={() => {
            // history.push('/checkout');
            // dispatch(toggleCartHidden());
          }}
        >
          GO TO CHECKOUT
        </CustomButton>
      </div>
    </div>
  );
};

export default CartDropdown;
