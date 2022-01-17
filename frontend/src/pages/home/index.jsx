import React from 'react';
import { useSelector } from 'react-redux';
import Box from '@material-ui/core/Box';
import { makeStyles } from '@material-ui/core/styles';
// import ContactForm from './ContactForm';

const useStyles = makeStyles(() => ({
  home: {
    position: 'relative',
    padding: 20,
    height: '100vw',
    display: 'flex',
    flexDirection: 'column',
  },
  mainImg: {
    width: 'auto',
    height: 200,
  },
  img: {
    width: '90vw',
  },
  homePageText: {
    fontSize: '40px',
    fontFamily: 'Verdana',
    fontWeight: 'bold',
    color: '#bbccee',
    textAlign: 'center',
    textShadow: '2px 2px #123f4d',
  },
}));

function Home() {
  const classes = useStyles();
  const user = useSelector((state) => state.user.currentUser);
  return (
    <Box className={classes.home} component="div" m={1}>
      <div className={classes.background} />
      <div className={classes.homePageText}>
        Welcome to our bakery {user && user.first_name}!
      </div>
      <div className={classes.homePageText}>
        {' '}
        The home of delightful and yummy desserts prepared with love and
        attention!
      </div>
      {/* <ContactForm /> */}
    </Box>
  );
}

export default Home;
