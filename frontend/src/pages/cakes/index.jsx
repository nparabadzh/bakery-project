import React, { useEffect, useState } from 'react';
import axios from 'axios';
import Box from '@material-ui/core/Box';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  root: {
    position: 'relative',
    padding: 20,
    width: '100%',
    height: '100vw',
  },
}));

function Cakes() {
  const [cakes, setCakes] = useState([]);

  useEffect(() => {
    axios.get(`/baked-goods`).then((res) => {
      const cakes = res.data;
      setCakes(cakes);
    });
  }, []);

  const classes = useStyles();
  return (
    <Box className={classes.root} component="div" m={1}>
      Cakes Page
      <ul>
        {cakes.map((el) => {
          return (
            <li>
              {el.name} - {el.price} - {el.type}
            </li>
          );
        })}
      </ul>
    </Box>
  );
}

export default Cakes;
