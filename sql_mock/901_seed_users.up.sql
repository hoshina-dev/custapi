INSERT INTO users (email, name, password, organization_id, description, research_categories)
SELECT
    v.email,
    v.name,
    'hashed-password',
    o.id,
    v.description,
    v.research_categories::TEXT[]
FROM (VALUES
    ('alice@chula.ac.th',   'Alice Korawit',       'Chulalongkorn University',                   'Associate Professor of Computational Biology at Chulalongkorn University.',  '{Bioinformatics,Genomics,"Systems Biology"}'),
    ('bob@sut.ac.th',       'Bob Sarawat',         'Suranaree University of Technology',         'Head of the School of Electrical Engineering, SUT.',                         '{Embedded Systems,IoT,"Power Electronics"}'),
    ('charlie@chula.ac.th', 'Charlie Suttipong',   'Chulalongkorn University',                   'PhD candidate in Materials Science at Chulalongkorn University.',            '{Materials Science,Nanotechnology}'),
    ('diana@nus.edu.sg',    'Diana Tan',           'School of Computing, NUS',                   'Director of the NUS Center for Quantum Technologies.',                       '{Quantum Computing,"Quantum Cryptography"}'),
    ('eve@u-tokyo.ac.jp',   'Eve Nakamura',        'Tokyo University',                           'Professor of Artificial Intelligence, University of Tokyo.',                 '{"Machine Learning","Deep Learning",NLP}'),
    ('frank@titech.ac.jp',  'Frank Sato',          'Tokyo Institute of Technology',              'Associate Dean of Research at Tokyo Institute of Technology.',               '{Robotics,"Control Systems",MEMS}'),
    ('grace@ethz.ch',       'Grace Mueller',       'ETH Zurich',                                 'Postdoctoral researcher in theoretical computer science.',                   '{Algorithms,"Complexity Theory",Cryptography}'),
    ('hiroshi@ethz.ch',     'Hiroshi Tanaka',      'ETH Zurich',                                 'JSPS Fellow with a joint position at ETH Zurich and Tokyo Tech.',            '{Photonics,"Quantum Optics","Laser Physics"}'),
    ('ivan@mit.edu',        'Ivan Chen',           'Massachusetts Institute of Technology',      'Graduate researcher in the MIT Computer Science department.',                '{Distributed Systems,"Cloud Computing"}'),
    ('julia@mit.edu',       'Julia Park',          'Massachusetts Institute of Technology',      'Lab director and PI of the MIT Autonomous Systems Lab.',                     '{Autonomous Vehicles,"Computer Vision",SLAM}'),
    ('kai@technion.ac.il',  'Kai Levi',            'Technion - Israel Institute of Technology',  'Senior researcher in the Technion Cyber-Security center.',                   '{Cybersecurity,"Formal Verification","Network Security"}'),
    ('walter.white@ethz.ch', 'Walter White',       'ETH Zurich',                                 'Lead chemist and research director at ETH Zurich.',                          '{Chemistry,"Materials Science","Chaos Theory"}')
) AS v(email, name, org_name, description, research_categories)
JOIN organizations o ON o.name = v.org_name AND o.deleted_at IS NULL
ON CONFLICT (email) DO NOTHING;
