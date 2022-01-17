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
  const [error, setError] = useState(false);

  const navigate = useNavigate();

  const dispatch = useDispatch();

  const signIn = () => {
    axios
      .post(`/login`, {
        email,
        password,
      })
      .then((res) => {
        if (res.data.message === 'logged in') {
          dispatch(setCurrentUser(res.data.user));
          setError(false);
          navigate('/');
        }
      })
      .catch((error) => {
        setError(true);
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
          {error && (
            <span style={{ color: 'red' }}>
              The email or password is incorrect!
            </span>
          )}
          <FormControl>
            <Button
              style={{ marginTop: 20 }}
              disabled={email === '' || password === ''}
              variant="contained"
              color="secondary"
              onClick={signIn}
            >
              Sign in
            </Button>
          </FormControl>
        </form>
      </div>
    </div>
  );
};

export default SignUp;
