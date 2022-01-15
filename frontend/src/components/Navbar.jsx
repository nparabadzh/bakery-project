import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
// import { createStructuredSelector } from 'reselect';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import MenuItem from '@material-ui/core/MenuItem';
import Menu from '@material-ui/core/Menu';
import AccountCircle from '@material-ui/icons/AccountCircle';
import CakeIcon from '@material-ui/icons/Cake';
import ShoppingCartIcon from '@material-ui/icons/ShoppingCart';

import { setCurrentUser } from '../redux/user/actions';
import { toggleCartHidden } from '../redux/cart/actions';
import { selectCartItemsCount } from '../redux/cart/selectors';
import CartDropdown from '../components/Cart/CartDropdown';

const useStyles = makeStyles((theme) => ({
  grow: {
    flexGrow: 1,
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    display: 'none',
    [theme.breakpoints.up('sm')]: {
      display: 'block',
    },
  },

  sectionDesktop: {
    display: 'none',
    [theme.breakpoints.up('md')]: {
      display: 'flex',
    },
  },
  navlinks: {
    marginLeft: theme.spacing(5),
  },
  link: {
    marginLeft: 20,
    textDecoration: 'none',
    color: 'white',
    fontSize: 20,
    '&:hover': {
      color: 'yellow',
      borderBottom: '1px solid white',
    },
  },
  cartIconDiv: {
    width: 45,
    height: 45,
    position: 'relative',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    cursor: 'pointer',
  },
}));

export default function PrimarySearchAppBar() {
  const classes = useStyles();
  const [anchorEl, setAnchorEl] = React.useState(null);
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const currentUser = useSelector((state) => state.user.currentUser);
  const hidden = useSelector((state) => state.cart.hidden);
  const itemCount = useSelector(selectCartItemsCount);

  const isMenuOpen = Boolean(anchorEl);

  const handleProfileMenuOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const signIn = () => {
    navigate('/signin');
  };

  const signUp = () => {
    navigate('/signup');
  };

  const signOut = () => {
    dispatch(setCurrentUser(null));
  };

  const toggleCart = () => {
    dispatch(toggleCartHidden());
  };

  const menuId = 'primary-search-account-menu';

  const renderSignInMenu = (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
      id={menuId}
      keepMounted
      transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      <MenuItem onClick={signIn}>Sign in</MenuItem>
      <MenuItem onClick={signUp}>Sign up</MenuItem>
    </Menu>
  );

  const renderSignOutMenu = (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
      id={menuId}
      keepMounted
      transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      <MenuItem onClick={signOut}>Sign out</MenuItem>
    </Menu>
  );

  return (
    <div className={classes.grow}>
      <AppBar style={{ background: '#2E3B55' }} position="static">
        <Toolbar>
          <Link to="/" className={classes.link}>
            <IconButton
              edge="start"
              className={classes.menuButton}
              color="inherit"
              aria-label="open drawer"
            >
              <CakeIcon />
              <Typography className={classes.title} variant="h6" noWrap>
                Naty&rsquo;s Cakes
              </Typography>
            </IconButton>
          </Link>
          <div className={classes.navlinks}>
            <Link to="/cakes" className={classes.link}>
              Boutique Cakes
            </Link>
            <Link to="/pies" className={classes.link}>
              Pies and tarts
            </Link>
            <Link to="/pastries" className={classes.link}>
              Pastries
            </Link>
            <Link to="/other" className={classes.link}>
              Other Goods
            </Link>
          </div>
          <div className={classes.grow} />
          <div className={classes.sectionDesktop}>
            {currentUser && (
              <div className={classes.cartIconDiv} onClick={toggleCart}>
                <ShoppingCartIcon />
                <span className={classes.itemCount}>{itemCount}</span>
              </div>
            )}

            {hidden ? null : <CartDropdown />}
            <IconButton
              edge="end"
              aria-label="account of current user"
              aria-controls={menuId}
              aria-haspopup="true"
              onClick={handleProfileMenuOpen}
              color="inherit"
            >
              <AccountCircle />
            </IconButton>
          </div>
        </Toolbar>
      </AppBar>
      {currentUser ? renderSignOutMenu : renderSignInMenu}
    </div>
  );
}
