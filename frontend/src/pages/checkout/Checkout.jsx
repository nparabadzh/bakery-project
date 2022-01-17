import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { makeStyles } from '@material-ui/core/styles';

import CheckoutItem from '../../components/CheckoutItem';

import { selectCartItems, selectCartTotal } from '../../redux/cart/selectors';
import CustomButton from '../../components/CustomButton';
import axios from 'axios';
import { useNavigate } from 'react-router';
import { clearCart } from '../../redux/cart/actions';

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
    fontSize: 25,
  },
}));

const pageSelector = createStructuredSelector({
  cartItems: selectCartItems,
  total: selectCartTotal,
});

const CheckoutPage = () => {
  const classes = useStyles();
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const { cartItems, total } = useSelector(pageSelector);
  const user = useSelector((state) => state.user.currentUser);

  const placeOrder = async () => {
    axios
      .post('/orders', {
        user_id: user.id,
        delivery_address: user.delivery_address,
        status: 'unprocessed',
      })
      .then((data) => {
        let orderId = data.data.id;
        cartItems.forEach((item) => {
          axios.post('/orderedGoods', {
            good_id: item.id,
            order_id: orderId,
            quantity: item.quantity,
          });
        });
        alert('Order is successfuly created!');
        dispatch(clearCart());
        navigate('/');
      })
      .catch((error) => console.log(error));
  };

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
      <div className={classes.total}>
        Total: {parseFloat(total).toFixed(2)} bgn
      </div>
      <CustomButton onClick={placeOrder} disabled={total === 0}>
        Place Order
      </CustomButton>
    </div>
  );
};

export default CheckoutPage;
