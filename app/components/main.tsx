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

function Main(): React.JSX.Element {
  const [isFetching, setIsFetching] = useState(false);
  // const isDarkMode = useColorScheme() === 'dark';
  const isDarkMode = false;

  const backgroundStyle = {
    backgroundColor: isDarkMode ? Colors.darker : Colors.lighter,
  };

  return (
    <View style={styles.buttonContainer}>
      <Button
        title="Delete Preferences"
        onPress={() => {
          deletePreferences(setIsFetching);
        }}
        disabled={isFetching}></Button>

      {isFetching && (
        <View>
          <ActivityIndicator
            size="small"
            color="white"
            style={{marginRight: 10}}
          />
          <Text style={styles.buttonText}>Deleting...</Text>
        </View>
      )}
    </View>
  );
}

async function deletePreferences(
  setIsFetching: React.Dispatch<React.SetStateAction<boolean>>,
) {
  setIsFetching(true);
  console.log('making api call');

  const timeout = 3000;
  const controller = new AbortController();
  const id = setTimeout(() => controller.abort(), timeout);

  try {
    const response = await fetch(
      'http://192.168.1.167:5000/api/delete-preference',
      {
        method: 'POST',
        signal: controller.signal,
      },
    );
    clearTimeout(id);
    const status = response.status;
    if (status != 200) {
      console.log('failed to fetch');
      ToastAndroid.show('Failed to delete your preference', ToastAndroid.SHORT);
    } else {
      console.log('successfully deleted');
      ToastAndroid.show(
        'Successfully deleted your preferences',
        ToastAndroid.SHORT,
      );
    }
    setIsFetching(false);
  } catch (error: any) {
    // console.error(error);
    if (error.name == 'AbortError') {
      ToastAndroid.show(
        'Request timeout, try after sometime',
        ToastAndroid.SHORT,
      );
    } else {
      ToastAndroid.show('Error deleting your preferences', ToastAndroid.SHORT);
    }
    setIsFetching(false);
  }
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

export default Main;
