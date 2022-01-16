import React from 'react';
import { useSelector } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { makeStyles } from '@material-ui/core/styles';

import CheckoutItem from '../../components/CheckoutItem';

import { selectCartItems, selectCartTotal } from '../../redux/cart/selectors';

const useStyles = makeStyles(() => ({
  checkoutPage: {
    padding: 20,
    background: 'white',
    width: '70%',
    hight: 'auto',
    minHeight: 300,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    margin: '50px auto 0',
  },
  checkoutHeader: {
    width: '100%',
    padding: '10px 0',
    display: 'flex',
    justifyContent: 'space-between',
    borderBottom: '1px solid darkgrey',
  },
  headerBlock: {
    textTransform: 'capitalize',
    width: '23%',
    '&:last-child': {
      width: '8%',
    },
  },
  total: {
    marginTop: 30,
    marginLeft: 'auto',
    fontSize: 36,
  },
}));

const pageSelector = createStructuredSelector({
  cartItems: selectCartItems,
  total: selectCartTotal,
});

const CheckoutPage = () => {
  const classes = useStyles();
  const { cartItems, total } = useSelector(pageSelector);
  return (
    <div className={classes.checkoutPage}>
      <div className={classes.checkoutHeader}>
        <div className={classes.headerBlock}>
          <span>Product</span>
        </div>
        <div className={classes.headerBlock}>
          <span>Description</span>
        </div>
        <div className={classes.headerBlock}>
          <span>Quantity</span>
        </div>
        <div className={classes.headerBlock}>
          <span>Price</span>
        </div>
        <div className={classes.headerBlock}>
          <span>Remove</span>
        </div>
      </div>
      {cartItems.map((cartItem) => (
        <CheckoutItem key={cartItem.id} cartItem={cartItem} />
      ))}
      <div className={classes.headerBlock}>TOTAL: ${total}</div>
    </div>
  );
};

export default CheckoutPage;
