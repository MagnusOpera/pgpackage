import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import styles from './index.module.css';

export default function Home() {
  return (
    <Layout
      title="PostgreSQL schema packaging for GitHub-native release flows"
      description="Build, diff, and apply PostgreSQL schema packages with a Go-native CLI."
    >
      <header className={styles.heroBanner}>
        <div className="container">
          <div className={styles.heroPanel}>
            <p className={styles.eyebrow}>Magnus Opera</p>
            <Heading as="h1" className={styles.heroTitle}>
              PostgreSQL schema packages built for reviewable releases
            </Heading>
            <p className={styles.heroSubtitle}>
              `pgpackage` turns SQL project sources into versioned artifacts, compares them with live targets,
              and applies changes with explicit destructive-operation gates.
            </p>
            <div className={styles.actions}>
              <Link className="button button--primary button--lg" to="/manual/learn/quickstart">
                Read the quickstart
              </Link>
              <Link className="button button--secondary button--lg" to="/manual/reference/project-file">
                Project file reference
              </Link>
            </div>
          </div>
        </div>
      </header>
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
