import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import styles from './index.module.css';

const principles = [
  {
    label: 'Desired state in SQL',
    text: 'Your schema files describe what PostgreSQL should look like now.',
  },
  {
    label: 'Git holds the history',
    text: 'Review normal SQL diffs instead of maintaining a migration narrative.',
  },
  {
    label: 'pgpac computes the delta',
    text: 'Build, compare, and apply against the live database only when you are ready.',
  },
];

const workflow = [
  {
    step: '01',
    title: 'Model the schema you want',
    text: 'Tables, views, functions, types, and security live as source-controlled SQL files.',
  },
  {
    step: '02',
    title: 'Diff desired state against reality',
    text: 'pgpac inspects the target database and generates an ordered plan from actual drift.',
  },
  {
    step: '03',
    title: 'Apply with reviewable safety gates',
    text: 'Teams can inspect text, JSON, or SQL output before execution, with destructive changes explicitly gated.',
  },
];

export default function Home() {
  return (
    <Layout
      title="DACPAC-style desired state for PostgreSQL"
      description="Treat PostgreSQL schema as desired state in Git, then let pgpac diff and apply it safely."
    >
      <header className={styles.heroBanner}>
        <div className="container">
          <div className={styles.heroPanel}>
            <Heading as="h1" className={styles.heroTitle}>
              PostgreSQL schema as desired state
            </Heading>
            <p className={styles.heroSubtitle}>
              `pgpac` is DACPAC for Postgres. Your SQL source is the desired state configuration, Git carries
              the change history, and `pgpac` computes the diff between that intent and a live database.
            </p>
            <div className={styles.actions}>
              <Link className="button button--primary button--lg" to="/manual/learn/quickstart">
                Read the quickstart
              </Link>
              <Link className="button button--secondary button--lg" to="/manual/">
                Why desired state
              </Link>
            </div>
            <div className={styles.principles}>
              {principles.map((item) => (
                <div key={item.label} className={styles.principleCard}>
                  <p className={styles.principleLabel}>{item.label}</p>
                  <p className={styles.principleText}>{item.text}</p>
                </div>
              ))}
            </div>
          </div>
          <div className={styles.codePanel}>
            <p className={styles.codeLabel}>Workflow</p>
            <pre className={styles.codeBlock}>
              <code>{`# SQL in Git is the source of truth
pgpac build --project app.pgpac --output out/

# Compare desired state to the live target
pgpac plan --package out/app.pgpkg --connection "postgres://..."

# Apply the reviewed delta
pgpac apply --package out/app.pgpkg --connection "postgres://..."`}</code>
            </pre>
            <p className={styles.codeFootnote}>
              Keep reasoning in terms of the schema you want. Let `pgpac` reason about the transition.
            </p>
          </div>
        </div>
      </header>
      <main>
        <section className={styles.storySection}>
          <div className="container">
            <div className={styles.storyIntro}>
              <Heading as="h2" className={styles.sectionTitle}>
                Stop encoding intent as migration choreography
              </Heading>
              <p className={styles.sectionText}>
                Migration-heavy workflows force teams to reason about intermediate steps, ordering, and drift repair.
                `pgpac` shifts the unit of work back to the schema itself: define the desired state in SQL, keep
                that state in Git, and plan from the target database&apos;s current reality.
              </p>
            </div>
            <div className={styles.workflowGrid}>
              {workflow.map((item) => (
                <div key={item.step} className={styles.workflowCard}>
                  <p className={styles.workflowStep}>{item.step}</p>
                  <Heading as="h3" className={styles.workflowTitle}>
                    {item.title}
                  </Heading>
                  <p className={styles.workflowText}>{item.text}</p>
                </div>
              ))}
            </div>
          </div>
        </section>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
