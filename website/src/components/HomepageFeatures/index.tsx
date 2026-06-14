import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: string;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Build portable schema packages',
    description:
      'Compile a PostgreSQL project file plus SQL sources into a self-contained .pgpkg artifact with manifest and checksums.',
  },
  {
    title: 'Inspect live drift before execution',
    description:
      'Generate ordered plans against a real target database and choose between text, JSON, or rendered SQL output.',
  },
  {
    title: 'Apply changes with explicit safety gates',
    description:
      'Destructive operations stay blocked until the project or operator explicitly opts in, keeping release automation predictable.',
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className={styles.featureCard}>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props) => (
            <Feature key={props.title} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
