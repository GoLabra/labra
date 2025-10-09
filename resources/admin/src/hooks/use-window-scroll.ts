// import { useEffect } from 'react';
// import { throttle } from 'lodash';

// interface Config {
//   handler: () => void;
//   delay?: number;
// }

// export const useWindowScroll = (config: Config): void => {
//   useEffect(
//     () => {
//       const { handler, delay } = config;

//       const withThrottle = throttle(handler, delay);

//       window.addEventListener('scroll', withThrottle);

//       return () => {
//         window.removeEventListener('scroll', withThrottle);
//       };
//     },
//     [config]
//   );
// };
