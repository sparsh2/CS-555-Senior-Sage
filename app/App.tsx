/**
 * Sample React Native App
 * https://github.com/facebook/react-native
 *
 * @format
 */

import React, {useState} from 'react';
import type {PropsWithChildren} from 'react';
import {
  SafeAreaView,
  ScrollView,
  StatusBar,
  StyleSheet,
  Text,
  useColorScheme,
  View,
  Button,
  ActivityIndicator,
  ToastAndroid,
} from 'react-native';

import {
  Colors,
  DebugInstructions,
  Header,
  LearnMoreLinks,
  ReloadInstructions,
} from 'react-native/Libraries/NewAppScreen';
import Main from './components/main';
import Login from './components/login';

type SectionProps = PropsWithChildren<{
  title: string;
}>;

function App(): React.JSX.Element {
  const [isFetching, setIsFetching] = useState(false);
  const [loggedIn, setLoggedIn] = useState(false);
  const isDarkMode = false;

  const backgroundStyle = {
    backgroundColor: isDarkMode ? Colors.darker : Colors.lighter,
  };
  
  //Callback function to handle login
  const handleLoginSuccess = () => {
    setLoggedIn(true); // This will display Main after login
  };

  //return <Main></Main>;
  return (
    <View style={{flex: 1}}>
      {loggedIn ? <Main /> : <Login onLoginSuccess={handleLoginSuccess} />}
    </View>
  );
}

const styles = StyleSheet.create({
  buttonContainer: {
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    flex: 1,
  },
  button: {
    flex: 1,
    padding: 10,
    backgroundColor: 'blue',
    borderRadius: 5,
  },
  buttonText: {
    color: 'white',
    fontSize: 16,
  },
});

export default App;
