-- Insert users (without organization_id, which was removed in migration 010)
INSERT INTO users (email, name, password, description, research_categories)
VALUES
    ('alice@chula.ac.th',    'Alice Korawit',     'hashed-password', 'Associate Professor of Computational Biology at Chulalongkorn University.',  '{Bioinformatics,Genomics,"Systems Biology"}'),
    ('bob@sut.ac.th',        'Bob Sarawat',       'hashed-password', 'Head of the School of Electrical Engineering, SUT.',                         '{Embedded Systems,IoT,"Power Electronics"}'),
    ('charlie@chula.ac.th',  'Charlie Suttipong', 'hashed-password', 'PhD candidate in Materials Science at Chulalongkorn University.',            '{Materials Science,Nanotechnology}'),
    ('diana@nus.edu.sg',     'Diana Tan',         'hashed-password', 'Director of the NUS Center for Quantum Technologies.',                       '{Quantum Computing,"Quantum Cryptography"}'),
    ('eve@u-tokyo.ac.jp',    'Eve Nakamura',      'hashed-password', 'Professor of Artificial Intelligence, University of Tokyo.',                 '{"Machine Learning","Deep Learning",NLP}'),
    ('frank@titech.ac.jp',   'Frank Sato',        'hashed-password', 'Associate Dean of Research at Tokyo Institute of Technology.',               '{Robotics,"Control Systems",MEMS}'),
    ('grace@ethz.ch',        'Grace Mueller',     'hashed-password', 'Postdoctoral researcher in theoretical computer science.',                   '{Algorithms,"Complexity Theory",Cryptography}'),
    ('hiroshi@ethz.ch',      'Hiroshi Tanaka',    'hashed-password', 'JSPS Fellow with a joint position at ETH Zurich and Tokyo Tech.',            '{Photonics,"Quantum Optics","Laser Physics"}'),
    ('ivan@mit.edu',         'Ivan Chen',         'hashed-password', 'Graduate researcher in the MIT Computer Science department.',                '{Distributed Systems,"Cloud Computing"}'),
    ('julia@mit.edu',        'Julia Park',        'hashed-password', 'Lab director and PI of the MIT Autonomous Systems Lab.',                     '{Autonomous Vehicles,"Computer Vision",SLAM}'),
    ('kai@technion.ac.il',   'Kai Levi',          'hashed-password', 'Senior researcher in the Technion Cyber-Security center.',                   '{Cybersecurity,"Formal Verification","Network Security"}'),
    ('walter.white@ethz.ch', 'Walter White',      'hashed-password', 'Lead chemist and research director at ETH Zurich.',                          '{Chemistry,"Materials Science","Chaos Theory"}')
ON CONFLICT (email) DO NOTHING;

-- Seed user_organizations (many-to-many with role)
INSERT INTO user_organizations (user_id, organization_id, role)
SELECT u.id, o.id, v.role
FROM (VALUES
    ('alice@chula.ac.th',    'Chulalongkorn University',                  'user'),
    ('bob@sut.ac.th',        'Suranaree University of Technology',        'user'),
    ('charlie@chula.ac.th',  'Chulalongkorn University',                  'user'),
    ('diana@nus.edu.sg',     'School of Computing, NUS',                  'admin'),
    ('eve@u-tokyo.ac.jp',    'Tokyo University',                          'user'),
    ('frank@titech.ac.jp',   'Tokyo Institute of Technology',             'admin'),
    ('grace@ethz.ch',        'ETH Zurich',                                'user'),
    ('hiroshi@ethz.ch',      'ETH Zurich',                                'user'),
    ('hiroshi@ethz.ch',      'Tokyo Institute of Technology',             'user'),
    ('ivan@mit.edu',         'Massachusetts Institute of Technology',     'user'),
    ('julia@mit.edu',        'Massachusetts Institute of Technology',     'admin'),
    ('kai@technion.ac.il',   'Technion - Israel Institute of Technology', 'user'),
    ('walter.white@ethz.ch', 'ETH Zurich',                                'admin')
) AS v(email, org_name, role)
JOIN users u ON u.email = v.email AND u.deleted_at IS NULL
JOIN organizations o ON o.name = v.org_name AND o.deleted_at IS NULL
ON CONFLICT (user_id, organization_id) DO NOTHING;
