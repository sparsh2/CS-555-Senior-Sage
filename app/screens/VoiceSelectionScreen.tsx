import React, { useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image } from 'react-native';
import Slider from '@react-native-community/slider';

export default function VoiceSelectionScreen({ navigation }: any) {
  const [selectedVoice, setSelectedVoice] = useState('Sol'); // Default voice
  const voices = [
    { name: 'Alloy', description: 'Energetic and bright' },
    { name: 'Echo', description: 'Calm and soothing' },
    { name: 'Fable', description: 'Creative and friendly' },
    { name: 'Onyx', description: 'Sharp and clear' },
    { name: 'Nova', description: 'Soft and gentle' },
    { name: 'Shimmer', description: 'Smooth and confident' },
  ];

  const handleVoiceChange = (index: number) => {
    setSelectedVoice(voices[index].name);
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>VA Voice Preference</Text>
      <Text style={styles.step}>Step 3/3</Text>
      <Text style={styles.chooseVoice}>Choose a voice</Text>

      <View style={styles.voiceContainer}>
        {voices.map((voice, index) => (
          <TouchableOpacity 
            key={voice.name} 
            style={[styles.voiceOption, selectedVoice === voice.name && styles.selectedVoice]} 
            onPress={() => handleVoiceChange(index)}
          >
            <Text style={styles.voiceName}>{voice.name}</Text>
            <Text style={styles.voiceDescription}>{voice.description}</Text>
          </TouchableOpacity>
        ))}
      </View>

      <Slider
        style={styles.slider}
        minimumValue={0}
        maximumValue={5}
        step={1}
        onValueChange={(value: number) => handleVoiceChange(value)}
        value={voices.findIndex((voice) => voice.name === selectedVoice)}
        minimumTrackTintColor="#FFFFFF"
        maximumTrackTintColor="#000000"
      />

      <Text style={styles.selectedVoiceText}>Selected Voice: {selectedVoice}</Text>
      

      <TouchableOpacity style={styles.skipButton} onPress={() => navigation.navigate('Home')}>
        <Text style={styles.skipText}>Skip</Text>
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#fff',
    padding: 20,
  },
  title: {
    fontSize: 24,
    color: '#000',
  },
  step: {
    fontSize: 16,
    color: '#000',
    marginTop: 10,
  },
  chooseVoice: {
    fontSize: 20,
    color: '#000',
    marginTop: 20,
  },
  voiceContainer: {
    marginTop: 30,
    width: '100%',
  },
  voiceOption: {
    backgroundColor: '#333',
    padding: 15,
    borderRadius: 10,
    marginBottom: 10,
    alignItems: 'center',
  },
  selectedVoice: {
    backgroundColor: '#555',
  },
  voiceName: {
    fontSize: 18,
    color: '#FFF',
  },
  voiceDescription: {
    fontSize: 14,
    color: '#AAA',
  },
  slider: {
    width: '80%',
    marginTop: 20,
  },
  selectedVoiceText: {
    fontSize: 18,
    color: '#FFF',
    marginTop: 10,
  },
  playButton: {
    marginTop: 30,
  },
  playIcon: {
    width: 50,
    height: 50,
  },
  skipButton: {
    marginTop: 5,
    padding: 10,
    backgroundColor: '#555',
    borderRadius: 20,
  },
  skipText: {
    color: '#FFF',
    fontSize: 16,
  },
});