import React from 'react';

interface Props {
  className?: string;
}

export const Test: React.FC<Props> = ({ className }) => {
  return <div className={className}></div>;
};
