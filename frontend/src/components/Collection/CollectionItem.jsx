import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import CustomButton from '../CustomButton';
import { useDispatch } from 'react-redux';
import { addItem } from '../../redux/cart/actions';

const useStyles = makeStyles(() => ({
  collectionItem: {
    widnt: '22cw',
    display: 'flex',
    flexDirection: 'column',
    height: 350,
    alignItems: 'center',
    position: 'relative',
    marginBottom: 30,
    backgroundColor: 'white',
    color: 'black',
    // border: '1px solid black',
  },
  image: {
    width: '90%',
    height: '90%',
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    marginBlock: 5,
  },
  footer: {
    width: '100%',
    height: '10%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-around',
    fontSize: 18,
    margin: 10,
  },
  name: {
    textAlign: 'center',
    marginBottom: 15,
    marginTop: 15,
  },
  price: {
    textAlign: 'center',
    marginBottom: 15,
  },
}));

const CollectionItem = ({ item, disabled }) => {
  const classes = useStyles();
  const dispatch = useDispatch();
  const { name, price, photo } = item;

  return (
    <div className={classes.collectionItem}>
      <div
        className={classes.image}
        style={{
          backgroundImage: `url(${photo})`,
        }}
      />
      <div className={classes.footer}>
        <span className={classes.name}>{name}</span>
        <span className={classes.price}>Price: {price} bgn</span>
      </div>
      <div style={{ padding: 10, margin: 10 }}>
        <CustomButton
          disabled={disabled}
          onClick={() => {
            dispatch(addItem(item));
          }}
        >
          Add to cart
        </CustomButton>
      </div>
    </div>
  );
};

export default CollectionItem;
