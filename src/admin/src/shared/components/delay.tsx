import React, { useState, useEffect, FC, ReactNode } from 'react';

type Props = {
	children: React.ReactNode;
	wait?: number;
	skeletonComponenet?: ReactNode;
};

const Delayed = ({ children, wait = 0, skeletonComponenet }: Props) => {
	const [isShown, setIsShown] = useState(false);

	useEffect(() => {
		const timer = setTimeout(() => {
			setIsShown(true);
		}, wait);
		return () => clearTimeout(timer);
	}, [wait]);

	return isShown ? children : skeletonComponenet;
};

export default Delayed;