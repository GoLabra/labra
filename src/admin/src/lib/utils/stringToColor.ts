
export function stringToColor(string: string) {
    let hash = 0;
    let i;

    /* eslint-disable no-bitwise */
    for (i = 0; i < string.length; i += 1) {
        hash = string.charCodeAt(i) + ((hash << 5) - hash);
    }

    let color = '#';

    for (i = 0; i < 3; i += 1) {
        const value = (hash >> (i * 8)) & 0xff;
        color += `00${value.toString(16)}`.slice(-2);
    }
    /* eslint-enable no-bitwise */

    return color + '30';
}

export function getInitials(name: string) {
    return name.split(' ').slice(0, 2).map(word => word[0]).join('').trim().toUpperCase();
}

// export function stringAvatar(name: string) {
//     return {
//         sx: {
//             position: 'relative',
//             '&::before': {
//                                 content: '""',
//                                 position: 'absolute',
//                                 top: 0,
//                                 left: 0,
//                                 right: 0,
//                                 bottom: 0,
//                                 backgroundColor: stringToColor(name),
//                         }
//         },
//         children: name.split(' ').slice(0, 2).map(word => word[0]).join('').toUpperCase(),
//     };
// }