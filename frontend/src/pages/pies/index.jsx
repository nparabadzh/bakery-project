import React, { useEffect, useState } from 'react';
import { useSelector } from 'react-redux';
import axios from 'axios';
import Box from '@material-ui/core/Box';
import { makeStyles } from '@material-ui/core/styles';
import CollectionPage from '../../components/Collection/CollectionPage';

const useStyles = makeStyles(() => ({
  root: {
    position: 'relative',
    padding: 20,
    width: '100%',
    height: '100vw',
  },
}));

const CATEGORY_NAME = 'Pies';

function Pies() {
  const [cakes, setCakes] = useState([]);
  const allCategories = useSelector((state) => state.categories.categories);
  const cakeCategory = Object.keys(allCategories).find(
    (key) => allCategories[key] === CATEGORY_NAME,
  );

  useEffect(() => {
    if (cakeCategory) {
      axios
        .get(`/baked-goods`, { params: { category_id: cakeCategory } })
        .then((res) => {
          const cakes = res.data;
          setCakes(cakes);
        });
    }
  }, [cakeCategory]);

  const classes = useStyles();
  return (
    <Box className={classes.root} component="div" m={1}>
      {cakes.length > 0 && (
        <CollectionPage items={cakes} title={CATEGORY_NAME} />
      )}
    </Box>
  );
}

export default Pies;
