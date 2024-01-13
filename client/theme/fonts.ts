import { useFonts } from "expo-font";
import { StyleSheet } from 'react-native'

export const fonts = StyleSheet.create({
    Alata: {
        fontFamily: 'Alata'
    }
})

export const useThemeFonts = () => {
    const [fontsLoaded] = useFonts({
        'Alata': require('../assets/fonts/Alata-Regular.ttf'),
      });
      
      if (!fontsLoaded) {
        return false;
    }
    return true
}