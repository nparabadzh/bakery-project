import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import CollectionItem from './CollectionItem';

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
    gridTemplateColumns: '1fr 1fr 1fr 1fr',
    gridGap: 10,
  },
}));

const CollectionPage = ({ items, title }) => {
  const classes = useStyles();
  console.log(items);
  return (
    <div className={classes.page}>
      <h2 className={classes.title}>{title}</h2>
      <div className={classes.items}>
        {items.map((item) => (
          <CollectionItem key={item.id} item={item} />
        ))}
      </div>
    </div>
  );
};

export default CollectionPage;
