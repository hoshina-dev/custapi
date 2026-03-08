-- Insert users (without organization_id, which was removed in migration 010)
WITH pwd AS (
  SELECT '$2a$10$b8HxaDpATeyIhV5QPj.wveEcqbKbD4Netq8/w039NF5ZaclMr8.Ge' AS hash
)
INSERT INTO users (email, name, password, description, research_categories)
SELECT
    v.email,
    v.name,
    pwd.hash,
    v.description,
    v.research_categories::TEXT[]
FROM pwd, (VALUES
    ('alice@chula.ac.th',    'Alice Korawit',     'Associate Professor of Computational Biology at Chulalongkorn University.',  '{Bioinformatics,Genomics,"Systems Biology"}'),
    ('bob@sut.ac.th',        'Bob Sarawat',       'Head of the School of Electrical Engineering, SUT.',                         '{Embedded Systems,IoT,"Power Electronics"}'),
    ('charlie@chula.ac.th',  'Charlie Suttipong', 'PhD candidate in Materials Science at Chulalongkorn University.',            '{Materials Science,Nanotechnology}'),
    ('diana@nus.edu.sg',     'Diana Tan',         'Director of the NUS Center for Quantum Technologies.',                       '{Quantum Computing,"Quantum Cryptography"}'),
    ('eve@u-tokyo.ac.jp',    'Eve Nakamura',      'Professor of Artificial Intelligence, University of Tokyo.',                 '{"Machine Learning","Deep Learning",NLP}'),
    ('frank@titech.ac.jp',   'Frank Sato',        'Associate Dean of Research at Tokyo Institute of Technology.',               '{Robotics,"Control Systems",MEMS}'),
    ('grace@ethz.ch',        'Grace Mueller',     'Postdoctoral researcher in theoretical computer science.',                   '{Algorithms,"Complexity Theory",Cryptography}'),
    ('hiroshi@ethz.ch',      'Hiroshi Tanaka',    'JSPS Fellow with a joint position at ETH Zurich and Tokyo Tech.',            '{Photonics,"Quantum Optics","Laser Physics"}'),
    ('ivan@mit.edu',         'Ivan Chen',         'Graduate researcher in the MIT Computer Science department.',                '{Distributed Systems,"Cloud Computing"}'),
    ('julia@mit.edu',        'Julia Park',        'Lab director and PI of the MIT Autonomous Systems Lab.',                     '{Autonomous Vehicles,"Computer Vision",SLAM}'),
    ('kai@technion.ac.il',   'Kai Levi',          'Senior researcher in the Technion Cyber-Security center.',                   '{Cybersecurity,"Formal Verification","Network Security"}'),
    ('walter.white@ethz.ch', 'Walter White',      'Lead chemist and research director at ETH Zurich.',                          '{Chemistry,"Materials Science","Chaos Theory"}')
) AS v(email, name, description, research_categories)
ON CONFLICT (email) DO NOTHING;

-- Seed user_organizations (many-to-many with role)
-- Ensures each organization has at least one admin
INSERT INTO user_organizations (user_id, organization_id, role)
SELECT u.id, o.id, v.role
FROM (VALUES
    ('alice@chula.ac.th',    'Chulalongkorn University',                  'admin'),
    ('charlie@chula.ac.th',  'Chulalongkorn University',                  'user'),
    ('bob@sut.ac.th',        'Suranaree University of Technology',        'admin'),
    ('diana@nus.edu.sg',     'School of Computing, NUS',                  'admin'),
    ('eve@u-tokyo.ac.jp',    'Tokyo University',                          'admin'),
    ('frank@titech.ac.jp',   'Tokyo Institute of Technology',             'admin'),
    ('hiroshi@ethz.ch',      'Tokyo Institute of Technology',             'user'),
    ('grace@ethz.ch',        'ETH Zurich',                                'user'),
    ('hiroshi@ethz.ch',      'ETH Zurich',                                'user'),
    ('walter.white@ethz.ch', 'ETH Zurich',                                'admin'),
    ('ivan@mit.edu',         'Massachusetts Institute of Technology',     'user'),
    ('julia@mit.edu',        'Massachusetts Institute of Technology',     'admin'),
    ('kai@technion.ac.il',   'Technion - Israel Institute of Technology', 'admin')
) AS v(email, org_name, role)
JOIN users u ON u.email = v.email AND u.deleted_at IS NULL
JOIN organizations o ON o.name = v.org_name AND o.deleted_at IS NULL
ON CONFLICT (user_id, organization_id) DO NOTHING;
