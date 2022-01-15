import React from 'react';
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

function Pies() {
  const classes = useStyles();
  return (
    <Box className={classes.root} component="div" m={1}>
      Pies Page
    </Box>
  );
}

export default Pies;
