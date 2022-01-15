import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import Button from '@material-ui/core/Button';
import FormControl from '@material-ui/core/FormControl';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import { setCurrentUser } from '../../redux/user/actions';

const SignUp = () => {
  const [email, setEmail] = useState('');
  const [password, setpassword] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [deliveryAddress, setDeliveryAddress] = useState('');

  const navigate = useNavigate();

  const dispatch = useDispatch();

  const signIn = () => {
    axios
      .post(`/signUp`, {
        email,
        password,
        first_name: firstName,
        last_name: lastName,
        delivery_address: deliveryAddress,
      })
      .then((res) => {
        debugger;
        if (res.statusText === 'Created') {
          dispatch(setCurrentUser(res.data));
          navigate('/');
        }
      })
      .catch((error) => {
        console.log(error);
      });
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center' }}>
      <div style={{ width: '60%', margin: 30 }}>
        <form
          style={{
            padding: 30,
            backgroundColor: 'white',
            display: 'flex',
            flexDirection: 'column',
          }}
          noValidate
          autoComplete="off"
        >
          <FormControl>
            <InputLabel htmlFor="component-simple">Email</InputLabel>
            <Input
              value={email}
              onChange={(e) => {
                setEmail(e.target.value);
              }}
            />
          </FormControl>
          <FormControl>
            <InputLabel htmlFor="component-helper">Password</InputLabel>
            <Input
              type="password"
              value={password}
              onChange={(e) => {
                setpassword(e.target.value);
              }}
              aria-describedby="component-helper-text"
            />
          </FormControl>
          <FormControl>
            <InputLabel htmlFor="component-simple">First Name</InputLabel>
            <Input
              value={firstName}
              onChange={(e) => {
                setFirstName(e.target.value);
              }}
            />
          </FormControl>
          <FormControl>
            <InputLabel htmlFor="component-helper">LastName</InputLabel>
            <Input
              value={lastName}
              onChange={(e) => {
                setLastName(e.target.value);
              }}
              aria-describedby="component-helper-text"
            />
          </FormControl>
          <FormControl>
            <InputLabel htmlFor="component-helper">Delivery address</InputLabel>
            <Input
              value={deliveryAddress}
              onChange={(e) => {
                setDeliveryAddress(e.target.value);
              }}
              aria-describedby="component-helper-text"
            />
          </FormControl>
          <FormControl>
            <Button
              style={{ marginTop: 20 }}
              disabled={
                email === '' ||
                password === '' ||
                firstName === '' ||
                lastName === '' ||
                deliveryAddress == ''
              }
              variant="contained"
              color="primary"
              onClick={signIn}
            >
              Sign up
            </Button>
          </FormControl>
        </form>
      </div>
    </div>
  );
};

export default SignUp;
