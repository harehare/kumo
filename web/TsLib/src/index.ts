import * as firebase from "firebase/app";
import "firebase/auth";

interface User {
  id: string;
  displayName: string;
  photoUrl: string;
  idToken: string;
}

const firebaseConfig = {
  apiKey: process.env.FIREBASE_API_KEY,
  authDomain: process.env.FIREBASE_AUTH_DOMAIN,
  projectId: process.env.FIREBASE_PROJECT_ID,
  appId: process.env.FIREBASE_APP_ID
};
firebase.initializeApp(firebaseConfig);

export const signIn = async _ => {
  const user = await firebase.auth().signInAnonymously();
  const idToken = await firebase
    .auth()
    .currentUser.getIdToken(true)
    .catch(() => {
      // Handle error
    });
  if (idToken) {
    markUserAsAuthenticated({
      id: user.user.uid,
      displayName: user.user.displayName,
      photoUrl: user.user.photoURL,
      idToken
    });
  }
};

export const signOut = async _ => {
  await firebase.auth().signOut();
  markUserAsAnonymous();
};

const markUserAsAuthenticated = (user: User) => {
  // @ts-ignore
  DotNet.invokeMethod("Kumo", "MarkUserAsAuthenticated", JSON.stringify(user));
};

const markUserAsAnonymous = () => {
  // @ts-ignore
  DotNet.invokeMethod("Kumo", "MarkUserAsAnonymous");
};

// @ts-ignore
Blazor.start({}).then(() => {
  firebase.auth().onAuthStateChanged(user => {
    if (user) {
      user.getIdToken(true).then(idToken => {
        markUserAsAuthenticated({
          id: user.uid,
          displayName: user.displayName,
          photoUrl: user.photoURL,
          idToken
        });
      });
    } else {
      markUserAsAnonymous();
    }
  });
});
