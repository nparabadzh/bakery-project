import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import CollectionItem from './CollectionItem';
import { useSelector } from 'react-redux';

const useStyles = makeStyles(() => ({
  page: {
    display: 'flex',
    flexDirection: 'column',
  },
  title: {
    fontSize: 38,
    margin: '0 auto 30px',
  },
  items: {
    display: 'grid',
    width: '95%',
    gridTemplateColumns: '1fr 1fr 1fr 1fr',
    gridGap: 10,
  },
}));

const CollectionPage = ({ items, title }) => {
  const classes = useStyles();
  const currentUser = useSelector((state) => state.user.currentUser);
  return (
    <div className={classes.page}>
      <h2 className={classes.title}>{title}</h2>
      <div className={classes.items}>
        {items.map((item) => (
          <CollectionItem
            key={item.id}
            item={item}
            disabled={currentUser == null}
          />
        ))}
      </div>
    </div>
  );
};

export default CollectionPage;
