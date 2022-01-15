import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles(() => ({
  customButton: {
    minWidth: '165px',
    width: 'auto',
    height: '50px',
    letterSpacing: '0.5px',
    lineHeight: '50px',
    padding: '0 35px 0 35px',
    fontSize: '15px',
    backgroundColor: 'black',
    color: 'white',
    textTransform: 'uppercase',
    fontFamily: "'Open Sans Condensed'",
    fontWeight: 'bolder',
    border: 'none',
    cursor: 'pointer',
    display: 'flex',
    justifyContent: 'center',
    '&:hover': {
      backgroundColor: 'white',
      color: 'black',
      border: '1px solid black',
    },
  },
}));

const CustomButton = ({ children, ...otherProps }) => {
  const classes = useStyles();
  return (
    <button className={classes.customButton} {...otherProps}>
      {children}
    </button>
  );
};

export default CustomButton;
