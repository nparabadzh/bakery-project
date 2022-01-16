import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router';

import CustomButton from '../CustomButton';
import CartItem from './CartItem';
import { toggleCartHidden } from '../../redux/cart/actions';

const useStyles = makeStyles(() => ({
  cartDropdown: {
    position: 'absolute',
    width: 240,
    height: 340,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    padding: 20,
    border: '1px solid black',
    backgroundColor: 'white',
    top: 50,
    right: 10,
    zIndex: 5,
  },
  emptyMsg: {
    fontSize: 18,
    color: 'black',
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

const CartDropdown = () => {
  const classes = useStyles();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const cartItems = useSelector((state) => state.cart.cartItems);
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
          disabled={cartItems.length === 0}
          onClick={() => {
            navigate('/checkout');
            dispatch(toggleCartHidden());
          }}
        >
          GO TO CHECKOUT
        </CustomButton>
      </div>
    </div>
  );
};

export default CartDropdown;
