import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import CustomButton from '../CustomButton';

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
    border: '1px solid black',
  },
  image: {
    width: '100%',
    height: '95%',
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    marginBlock: 5,
  },
  footer: {
    width: '100%',
    height: '5%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-around',
    fontSize: 18,
    margin: 10,
  },
  name: {
    textAlign: 'center',
    marginBottom: 15,
  },
  price: {
    textAlign: 'center',
    marginBottom: 15,
  },
}));

const CollectionItem = ({ item }) => {
  const classes = useStyles();
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
      <div style={{ paddign: 10, margin: 10 }}>
        <CustomButton onClick={() => {}}>Add to cart</CustomButton>
      </div>
    </div>
  );
};

export default CollectionItem;
